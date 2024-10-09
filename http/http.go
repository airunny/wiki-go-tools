package grpc

import (
	"context"
	"os"
	"time"

	k8s "github.com/airunny/wiki-go-tools/kubernetes"
	mmd "github.com/airunny/wiki-go-tools/metadata"
	"github.com/airunny/wiki-go-tools/registry"
	"github.com/go-kratos/kratos/v2/log" // nolint
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

const (
	kubeConfigName        = "KUBE_CONFIG"
	defaultKubeConfigPath = "/kube/config"
)

type Option struct {
	timeout time.Duration
	logger  log.Logger
}

func newOption() *Option {
	return &Option{
		timeout: time.Second * 5,
		logger:  log.GetLogger(),
	}
}

type ClientOption func(o *Option)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *Option) {
		o.timeout = timeout
	}
}

func WithLogger(logger log.Logger) ClientOption {
	return func(o *Option) {
		o.logger = logger
	}
}

func NewClientWithShort(ctx context.Context, endpoint string, opts ...ClientOption) *kratosHttp.Client {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	clientOpts := []kratosHttp.ClientOption{
		kratosHttp.WithEndpoint(endpoint),
		kratosHttp.WithTimeout(o.timeout),
		kratosHttp.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
			tracing.Client(),
			mmd.Client(),
		),
	}

	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if ok && os.Getenv("REGISTRY") != "NO" {
		clientSet, err := k8s.NewClient()
		if err != nil {
			panic(err)
		}

		reg := registry.NewRegistry(clientSet, o.logger)
		reg.Start()

		clientOpts = append(clientOpts, kratosHttp.WithDiscovery(reg))
	}

	client, err := kratosHttp.NewClient(ctx, clientOpts...)
	if err != nil {
		panic(err)
	}
	return client
}

func NewClient(ctx context.Context, logger log.Logger, opts ...kratosHttp.ClientOption) (*kratosHttp.Client, error) {
	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if ok && os.Getenv("REGISTRY") != "NO" {
		clientSet, err := k8s.NewClient()
		if err != nil {
			panic(err)
		}

		reg := registry.NewRegistry(clientSet, logger)
		reg.Start()

		opts = append(opts, kratosHttp.WithDiscovery(reg))
	}

	return kratosHttp.NewClient(ctx, opts...)
}
