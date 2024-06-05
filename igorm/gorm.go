package igorm

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	dmysql "github.com/go-sql-driver/mysql"
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
	sourceFormat       = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local%v"
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

	DB struct {
		orm *gorm.DB
		db  *sql.DB
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

func NewGORM(c *Config, log log.Logger) (*DB, error) {
	var (
		err    error
		source = c.Source
	)

	if source == "" {
		useSsl := ""
		if len(c.CaPath) > 0 {
			var sslCfg *tls.Config
			sslCfg, err = GenTLSConfig(c.CaPath)
			if err != nil {
				return nil, err
			}

			_ = dmysql.RegisterTLSConfig("custom", sslCfg)
			useSsl = "&tls=custom"
		}
		source = fmt.Sprintf(sourceFormat, c.User, c.Password, c.Domain, c.Port, c.Database, useSsl)
	}

	if source == "" {
		return nil, fmt.Errorf("empty db source")
	}

	orm, err := gorm.Open(mysql.Open(source), &gorm.Config{
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
		return nil, err
	}

	err = orm.Callback().Create().Before("gorm:create").Register("create", Before)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Delete().Before("gorm:delete").Register("delete", Before)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Update().Before("gorm:update").Register("update", Before)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Query().Before("gorm:query").Register("query", Before)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Row().Before("gorm:row").Register("row", Before)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Create().After("gorm:create").Register("create_after", After)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Delete().After("gorm:delete").Register("delete_after", After)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Delete().After("gorm:update").Register("update_after", After)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Query().After("gorm:query").Register("query_after", After)
	if err != nil {
		return nil, err
	}
	err = orm.Callback().Row().After("gorm:row").Register("raw_after", After)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ormDB, err := orm.DB()
	if err != nil {
		return nil, err
	}

	err = ormDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	if c.MaxOpen <= 0 {
		c.MaxOpen = 100
	}

	if c.MaxIdle <= 0 {
		c.MaxIdle = 10
	}

	ormDB.SetMaxOpenConns(c.MaxOpen)
	ormDB.SetMaxIdleConns(c.MaxIdle)
	ormDB.SetConnMaxLifetime(c.MaxLifeTime)

	globalDB = ormDB
	globalGORM = orm
	return &DB{
		orm: orm,
		db:  ormDB,
	}, nil
}

func (d *DB) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *DB) Session() *gorm.DB {
	return d.orm
}

func (d *DB) NewOptions(opts ...Option) *Options {
	o := &Options{
		tx: d.orm,
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (d *DB) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return d.orm.Begin(opts...)
}

func GenTLSConfig(caCertFile string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}, nil
}

func typeFromModel(model interface{}) reflect.Type {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
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
