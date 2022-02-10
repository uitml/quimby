package k8s

import (
	"context"

	"github.com/uitml/quimby/internal/validate"
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
		TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}

	return &ns
}

func (c *Client) UserExists(u string) (bool, error) {
	namespaces, err := c.GetNamespaceList()
	if err != nil {
		return false, err
	}

	for _, namespace := range namespaces.Items {
		if validate.Username(namespace.Name) && namespace.Name == u {
			return true, nil
		}
	}

	return false, nil
}

func (c *Client) DeleteUser(u string) error {
	gracePeriod := int64(0)
	policy := metav1.DeletePropagationForeground
	opts := metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod, PropagationPolicy: &policy}

	err := c.Clientset.CoreV1().Namespaces().Delete(context.TODO(), u, opts)
	if err != nil {
		return err
	}
	return nil
}
