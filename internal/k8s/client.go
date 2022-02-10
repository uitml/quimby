package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceClient interface {
	GetNamespaceList() (*corev1.NamespaceList, error)
	GetResourceQuota(string) (ResourceQuota, error)
	GetDefaultRequest(string) (ResourceRequest, error)
	NewUser(string, string, string, string) error
	NewSimpleUser(string) error
	GetTotalGPUs() (ResourceSummary, error)
	UserExists(string) (bool, error)
	DeleteUser(string) error
}

type Client struct {
	Clientset kubernetes.Interface
}

func NewClient() (ResourceClient, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return &Client{kubernetes.NewForConfigOrDie(config)}, nil
}
