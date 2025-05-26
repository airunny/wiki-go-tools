package igorm

import (
	"context"
	"crypto/tls"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	opentracing "github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type dbDurationKey struct{}

const (
	traceOperationName = "GORM"
)

var (
	dbStatement = attribute.Key("db.statement")
	dbType      = attribute.Key("dbType")
	dbComponent = attribute.Key("component")
)

type (
	Config struct {
		Domain             string        `json:"domain" yaml:"domain"`
		Port               int           `json:"port" yaml:"port"`
		Database           string        `json:"database" yaml:"database"`
		User               string        `json:"user" yaml:"user" env:"DP_MYSQL_DB_USER"`
		Password           string        `json:"password" yaml:"password" env:"DP_MYSQL_DB_PASSWORD"`
		LogLevel           int           `json:"log_level" yaml:"logLevel" env:"DP_MYSQL_DB_LOG"`
		MaxOpen            int           `json:"max_open" yaml:"maxOpen" env:"DP_MYSQL_DB_MAX_OPEN" default:"100"`
		MaxIdle            int           `json:"max_idle" yaml:"maxIdle" env:"DP_MYSQL_DB_MAX_IDLE" default:"10"`
		MaxLifeTime        time.Duration `json:"max_life_time" yaml:"maxLifeTime" env:"DP_MYSQL_DB_MAX_LIFE_TIME" default:"8"`
		DisableQueryFields bool          `json:"disable_query_fields" yaml:"disableQueryFields" env:"DP_DISABLE_QUERY_FIELDS"`
		CaPath             string        `json:"ca_path" yaml:"caPath"`
		TLSConfig          *tls.Config   `json:"tls_config" yaml:"tlsConfig"`
		Source             string        `json:"source" yaml:"source"`
	}
)

func Before(db *gorm.DB) {
	db.Statement.Context = context.WithValue(db.Statement.Context, dbDurationKey{}, time.Now())
}

func After(db *gorm.DB) {
	ctx := db.Statement.Context
	startTime := ctx.Value(dbDurationKey{}).(time.Time)
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		tracer := otel.GetTracerProvider().Tracer("GORM")
		_, span := tracer.Start(ctx,
			traceOperationName+"_"+db.Statement.Table,
			oteltrace.WithSpanKind(oteltrace.SpanKindClient),
			oteltrace.WithTimestamp(startTime),
		)
		span.SetAttributes(dbStatement.String(db.Statement.SQL.String()),
			dbType.String("sql"),
			dbComponent.String("gorm.io/gorm"),
		)
		defer span.End()

		if db.Error == nil || errors.Is(db.Error, sql.ErrNoRows) || errors.Is(db.Error, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "")
			return
		}
		span.SetStatus(codes.Error, db.Error.Error())
		span.RecordError(db.Error)
		return
	}
	opentracing.StartSpan(traceOperationName, opentracing.StartTime(startTime),
		// https://github.com/opentracing-contrib/opentracing-specification-zh/blob/master/semantic_conventions.md
		opentracing.Tags{
			"db.statement": db.Statement.SQL.String(),
			"db.type":      "sql",
			"component":    "gorm.io/gorm",
		},
		opentracing.ChildOf(parentSpan.Context())).Finish()
}

func NewGORM(c *Config, log log.Logger) (*gorm.DB, io.Closer, error) {
	var (
		err    error
		source = c.Source
	)

	if source == "" {
		return nil, nil, fmt.Errorf("empty db source")
	}

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		Logger:                                   newLogger(log, gormlogger.LogLevel(c.LogLevel)),
		SkipDefaultTransaction:                   true,
		QueryFields:                              c.DisableQueryFields,
		PrepareStmt:                              false,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	err = db.Callback().Create().Before("gorm:create").Register("create", Before)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Delete().Before("gorm:delete").Register("delete", Before)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Update().Before("gorm:update").Register("update", Before)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Query().Before("gorm:query").Register("query", Before)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Row().Before("gorm:row").Register("row", Before)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Create().After("gorm:create").Register("create_after", After)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Delete().After("gorm:delete").Register("delete_after", After)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Delete().After("gorm:update").Register("update_after", After)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Query().After("gorm:query").Register("query_after", After)
	if err != nil {
		return nil, nil, err
	}
	err = db.Callback().Row().After("gorm:row").Register("raw_after", After)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sqldb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	err = sqldb.PingContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	if c.MaxOpen <= 0 {
		c.MaxOpen = 100
	}

	if c.MaxIdle <= 0 {
		c.MaxIdle = 10
	}

	sqldb.SetMaxOpenConns(c.MaxOpen)
	sqldb.SetMaxIdleConns(c.MaxIdle)
	sqldb.SetConnMaxLifetime(c.MaxLifeTime)

	return db, sqldb, nil
}

func typeFromModel(model interface{}) reflect.Type {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

type CustomValue[T any] struct {
	V T `json:"v"`
}

func (s *CustomValue[T]) Value() (driver.Value, error) {
	return GormCustomValue(s)
}

func (s *CustomValue[T]) Scan(value interface{}) error {
	return GormCustomScan(s, value)
}

func GormCustomValue(in interface{}) (driver.Value, error) {
	if in == nil {
		return "", nil
	}

	str, _ := json.Marshal(in)
	return string(str), nil
}

func GormCustomScan(target, value interface{}) error {
	if target == nil || value == nil {
		target = reflect.New(typeFromModel(target)).Interface()
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		bytes = []byte("{}")
	}

	if bytes == nil || len(bytes) <= 0 {
		bytes = []byte("{}")
	}

	return json.Unmarshal(bytes, target)
}
