package iredis

import (
	"context"
	"errors"
	"time"

	redis "github.com/go-redis/redis/v8"
	opentracing "github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace" // nolint
)

var redisCommandAttributeKey = attribute.Key("db.statement")

const spanName = "redis"

type _traceKey struct{}

func newHook(l *logger) *hook {
	return &hook{
		tracer: otel.GetTracerProvider().Tracer("Redis"),
		log:    l,
	}
}

type hook struct {
	tracer oteltrace.Tracer
	log    *logger
}

func (h hook) startSpan(ctx context.Context, cmds ...redis.Cmder) context.Context {
	ctx, span := h.tracer.Start(ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	)

	cmdStrs := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdStrs = append(cmdStrs, cmd.String())
	}
	span.SetAttributes(redisCommandAttributeKey.StringSlice(cmdStrs))
	span.SetAttributes()
	return ctx
}

func (h hook) endSpan(ctx context.Context, err error) {
	span := oteltrace.SpanFromContext(ctx)
	defer span.End()

	if err == nil || err == redis.Nil {
		span.SetStatus(codes.Ok, "")
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

func (h hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	err := cmd.Err()
	h.endSpan(ctx, err)
	h.log.Printf(ctx, "RedisCommand:%s", cmd.String())
	return nil
}

func (h hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return h.startSpan(context.WithValue(ctx, _traceKey{}, time.Now()), cmd), nil
}

func (h hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if len(cmds) == 0 {
		return ctx, nil
	}

	return h.startSpan(context.WithValue(ctx, _traceKey{}, time.Now()), cmds...), nil
}

func (h hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if len(cmds) == 0 {
		return nil
	}

	batchErrors := make([]error, 0, len(cmds))
	for _, cmd := range cmds {
		err := cmd.Err()
		if err == nil {
			continue
		}
		batchErrors = append(batchErrors, err)
	}
	h.endSpan(ctx, errors.Join(batchErrors...))

	start, ok := ctx.Value(_traceKey{}).(time.Time)
	if !ok {
		return nil
	}
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		return nil
	}
	opentracing.StartSpan("redis_pipeline", opentracing.StartTime(start),
		opentracing.Tags{
			"component":    "github.com/go-redis/redis/v8",
			"db.type":      "redis",
			"db.statement": cmds,
		},
		opentracing.ChildOf(parentSpan.Context())).Finish()
	return nil
}
