package trace

import (
	"errors"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger" // nolint
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type options struct {
	name       string
	namespace  string
	version    string
	instance   string
	attributes []attribute.KeyValue
	fraction   string
}

func newDefaultOptions() *options {
	return &options{
		name:      os.Getenv("SERVICE_NAME"),
		namespace: os.Getenv("SERVICE_NAMESPACE"),
		version:   os.Getenv("SERVICE_VERSION"),
		instance:  os.Getenv("SERVICE_INSTANCE"),
	}
}

type Option func(*options)

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

func WithInstance(instance string) Option {
	return func(o *options) {
		o.instance = instance
	}
}

func WithAttributes(attrs []attribute.KeyValue) Option {
	return func(o *options) {
		o.attributes = attrs
	}
}

func WithFraction(fraction string) Option {
	return func(o *options) {
		o.fraction = fraction
	}
}

func NewTrace(kind, url string, opts ...Option) (trace.TracerProvider, error) {
	if url == "" {
		return noop.NewTracerProvider(), nil
	}

	// 目前只支持jaeger，后续可以根据kind增加
	return NewJaegerTrace(url, opts...)
}

func NewJaegerTrace(url string, opts ...Option) (trace.TracerProvider, error) {
	if url == "" {
		return nil, errors.New("empty trace url")
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	o := newDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return newTrace(exporter, o)
}

func newTrace(exporter tracesdk.SpanExporter, o *options) (trace.TracerProvider, error) {
	var (
		attrs []attribute.KeyValue
	)

	if o.name != "" {
		attrs = append(attrs, semconv.ServiceNameKey.String(o.name))
	}

	if o.namespace != "" {
		attrs = append(attrs, semconv.ServiceNamespaceKey.String(o.namespace))
	}

	if o.version != "" {
		attrs = append(attrs, semconv.ServiceVersionKey.String(o.version))
	}

	if o.instance != "" {
		attrs = append(attrs, semconv.ServiceInstanceIDKey.String(o.instance))
	}

	if len(o.attributes) > 0 {
		attrs = append(attrs, o.attributes...)
	}

	providerOptions := []tracesdk.TracerProviderOption{tracesdk.WithBatcher(exporter)}
	if o.fraction != "" {
		providerOptions = append(providerOptions, tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))))
	}

	if len(attrs) > 0 {
		tracesdk.WithResource(resource.NewSchemaless(
			attrs...,
		))
	}

	tp := tracesdk.NewTracerProvider(providerOptions...)
	otel.SetTracerProvider(tp)
	return tp, nil
}
