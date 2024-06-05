package metrics

import (
	"context"
	"strconv"
	"time"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duration(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1, 2, 3, 4, 5},
	}, []string{"service", "namespace", "kind", "operation"})

	metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"service", "namespace", "kind", "operation", "code", "reason"})
)

func init() {
	prometheus.MustRegister(metricSeconds, metricRequests)
}

func ServerMetricsMiddleware(serviceName string) middleware.Middleware {
	return Server(
		serviceName,
		WithSeconds(prom.NewHistogram(metricSeconds)),
		WithRequests(prom.NewCounter(metricRequests)),
	)
}

func ClientMetricsMiddleware(serviceName string) middleware.Middleware {
	return Client(
		serviceName,
		WithSeconds(prom.NewHistogram(metricSeconds)),
		WithRequests(prom.NewCounter(metricRequests)),
	)
}

// Option is metrics option.
type Option func(*options)

// WithRequests with requests counter.
func WithRequests(c metrics.Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
func WithSeconds(c metrics.Observer) Option {
	return func(o *options) {
		o.seconds = c
	}
}

type options struct {
	requests metrics.Counter
	seconds  metrics.Observer
}

func Server(serviceName string, opts ...Option) middleware.Middleware {
	var (
		op        = options{}
		namespace = "server"
	)

	for _, o := range opts {
		o(&op)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if op.requests != nil {
				op.requests.With(serviceName, namespace, kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if op.seconds != nil {
				op.seconds.With(serviceName, namespace, kind, operation).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}

func Client(serviceName string, opts ...Option) middleware.Middleware {
	var (
		op        = options{}
		namespace = "client"
	)

	for _, o := range opts {
		o(&op)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if op.requests != nil {
				op.requests.With(serviceName, namespace, kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if op.seconds != nil {
				op.seconds.With(serviceName, namespace, kind, operation).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}
