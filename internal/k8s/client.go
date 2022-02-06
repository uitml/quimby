package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	GetNamespaceList() (*corev1.NamespaceList, error)
	GetResourceQuota(string) (ResourceQuota, error)
	GetDefaultRequest(string) (ResourceRequest, error)
	GetTotalGPUs() (int, error)
}

type K8sClient struct {
	Clientset *kubernetes.Clientset
}

func NewClient() (Client, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return &K8sClient{kubernetes.NewForConfigOrDie(config)}, nil
}

func (c *K8sClient) GetNamespaceList() (*corev1.NamespaceList, error) {
	namespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaceList, nil
}
