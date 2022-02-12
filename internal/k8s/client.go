package k8s

import (
	"github.com/uitml/quimby/internal/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceClient interface {
	NamespaceList() (*corev1.NamespaceList, error)
	Quota(string) (resource.Quota, error)
	DefaultRequest(string) (resource.Request, error)
	Namespace(string) (*corev1.Namespace, error)
	NewUser(string, string, string, string) error
	ApplyMetaData(string, string, string, string) error
	NewSimpleUser(string) error
	TotalGPUs() (resource.Summary, error)
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
