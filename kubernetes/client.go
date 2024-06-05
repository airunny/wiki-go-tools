package kubernetes

import (
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	kubeConfigName = "KUBE_CONFIG_PATH"
)

func NewClient() (*k8s.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return k8s.NewForConfig(config)
}
