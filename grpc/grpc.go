package grpc

import (
	"context"
	"os"

	k8s "github.com/airunny/wiki-go-tools/kubernetes"
	"github.com/airunny/wiki-go-tools/registry"
	"github.com/go-kratos/kratos/v2/log" // nolint
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

const (
	kubeConfigName        = "KUBE_CONFIG"
	defaultKubeConfigPath = "/kube/config"
)

func DialInsecure(ctx context.Context, logger log.Logger, opts ...kratosGrpc.ClientOption) (*grpc.ClientConn, error) {
	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if ok {
		clientSet, err := k8s.NewClient()
		if err != nil {
			panic(err)
		}

		reg := registry.NewRegistry(clientSet, logger)
		reg.Start()

		opts = append(opts, kratosGrpc.WithDiscovery(reg))
	}

	return kratosGrpc.DialInsecure(ctx, opts...)
}
