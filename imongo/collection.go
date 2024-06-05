package imongo

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	MongoName = "Mongo"
)

var (
	dbStatement = attribute.Key("db.statement")
	dbType      = attribute.Key("dbType")
	dbComponent = attribute.Key("component")
)

type Collection struct {
	*mongo.Collection
	debug bool
}

func getCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

func (c *Collection) final() func(ctx context.Context, end time.Time, kvs ...interface{}) {
	var (
		start    = time.Now()
		funcName = getCurrentFunctionName()
	)

	return func(ctx context.Context, end time.Time, kvs ...interface{}) {
		c.trace(ctx, start, funcName)
		c.logging(ctx, funcName, kvs...)
	}
}

func (c *Collection) trace(ctx context.Context, start time.Time, operation string) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		tracer := otel.GetTracerProvider().Tracer(MongoName)
		_, span := tracer.Start(ctx,
			fmt.Sprintf("%v:%v", c.Name(), operation),
			oteltrace.WithSpanKind(oteltrace.SpanKindClient),
			oteltrace.WithTimestamp(start),
		)
		span.SetAttributes(dbStatement.String(operation),
			dbType.String("sql"),
			dbComponent.String("mongo"),
		)
		defer span.End()
		return
	}

	opentracing.StartSpan(MongoName, opentracing.StartTime(start), opentracing.Tags{
		"action_name": operation,
		"component":   "go.mongodb.org/mongo-driver/mongo",
		"db.type":     "mongo",
	}, opentracing.ChildOf(parentSpan.Context())).Finish()
}

func (c *Collection) logging(ctx context.Context, operation string, kv ...interface{}) {
	if !c.debug {
		return
	}

	l := log.Context(ctx)
	keyValues := make([]interface{}, 0, len(kv)+2)
	keyValues = append(keyValues, MongoName, operation)
	keyValues = append(keyValues, kv...)
	l.Infow(keyValues...)
}

func (c *Collection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	defer c.final()(ctx, time.Now(), "Document", document)
	resp, err := c.Collection.InsertOne(ctx, document, opts...)
	return resp, err
}

func (c *Collection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	defer c.final()(ctx, time.Now(), "documents", documents)
	return c.Collection.InsertMany(ctx, documents, opts...)
}

func (c *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter)
	return c.Collection.Find(ctx, filter, opts...)
}

func (c *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	defer c.final()(ctx, time.Now(), "Filter", filter)
	return c.Collection.FindOne(ctx, filter, opts...)
}

func (c *Collection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	defer c.final()(ctx, time.Now(), "Filter", filter)
	return c.Collection.FindOneAndDelete(ctx, filter, opts...)
}

func (c *Collection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	defer c.final()(ctx, time.Now(), "Filter", filter, "Replacement", replacement)
	return c.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (c *Collection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	defer c.final()(ctx, time.Now(), "Filter", filter, "Update", update)
	return c.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter, "Update", update)
	return c.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter, "Update", update)
	return c.Collection.UpdateMany(ctx, filter, update, opts...)
}

func (c *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	defer c.final()(ctx, time.Now(), "Id", id, "Update", update)
	return c.Collection.UpdateByID(ctx, id, update, opts...)
}

func (c *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter)
	return c.Collection.DeleteOne(ctx, filter, opts...)
}

func (c *Collection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter)
	return c.Collection.DeleteMany(ctx, filter, opts...)
}

func (c *Collection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter)
	return c.Collection.CountDocuments(ctx, filter, opts...)
}

func (c *Collection) ReplaceOne(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	defer c.final()(ctx, time.Now(), "Filter", filter, "Replacement", replacement)
	return c.Collection.ReplaceOne(ctx, filter, replacement, opts...)
}

func (c *Collection) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	defer c.final()(ctx, time.Now(), "EstimatedDocumentCount", time.Now())
	return c.Collection.EstimatedDocumentCount(ctx, opts...)
}

func (c *Collection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]interface{}, error) {
	defer c.final()(ctx, time.Now(), "FieldName", fieldName, "Filter", filter)
	return c.Collection.Distinct(ctx, fieldName, filter, opts...)
}

func (c *Collection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	defer c.final()(ctx, time.Now(), "BulkWrite", time.Now())
	return c.Collection.BulkWrite(ctx, models, opts...)
}

func (c *Collection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	defer c.final()(ctx, time.Now(), "Aggregate", time.Now())
	return c.Collection.Aggregate(ctx, pipeline, opts...)
}
