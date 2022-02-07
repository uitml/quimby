package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AnnotationUserFullname string = "springfield.uit.no/user-fullname"
	AnnotationUserEmail    string = "springfield.uit.no/user-email"
	LabelUserType          string = "springfield.uit.no/user-type"
)

func (c *Client) GetNamespaceList() (*corev1.NamespaceList, error) {
	namespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaceList, nil
}

func NewNamespace(name string, labels map[string]string, annotations map[string]string) *corev1.Namespace {
	ns := corev1.Namespace{
		metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		corev1.NamespaceSpec{},
		corev1.NamespaceStatus{},
	}

	return &ns
}
