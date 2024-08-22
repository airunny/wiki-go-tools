package grpc

import (
	"context"
	"os"
	"time"

	k8s "github.com/airunny/wiki-go-tools/kubernetes"
	"github.com/airunny/wiki-go-tools/registry"
	"github.com/go-kratos/kratos/v2/log" // nolint
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
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

func DialInsecure(ctx context.Context, endpoint string, opts ...ClientOption) *grpc.ClientConn {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	clientOpts := []kratosGrpc.ClientOption{
		kratosGrpc.WithEndpoint(endpoint),
		kratosGrpc.WithTimeout(o.timeout),
		kratosGrpc.WithMiddleware(
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

		clientOpts = append(clientOpts, kratosGrpc.WithDiscovery(reg))
	}

	conn, err := kratosGrpc.DialInsecure(ctx, clientOpts...)
	if err != nil {
		panic(err)
	}
	return conn
}
