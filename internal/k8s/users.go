package k8s

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/uitml/quimby/internal/validate"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const (
	AnnotationUserFullname string = "springfield.uit.no/user-fullname"
	AnnotationUserEmail    string = "springfield.uit.no/user-email"
	LabelUserType          string = "springfield.uit.no/user-type"
)

func (c *Client) NamespaceList() (*corev1.NamespaceList, error) {
	namespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaceList, nil
}

func (c *Client) Namespace(username string) (*corev1.Namespace, error) {
	namespaces, err := c.NamespaceList()
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces.Items {
		if validate.Username(namespace.Name) && namespace.Name == username {
			return &namespace, nil
		}
	}

	return nil, fmt.Errorf("user %s does not exist", username)

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
	namespaces, err := c.NamespaceList()
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
	policy := metav1.DeletePropagationForeground
	opts := metav1.DeleteOptions{GracePeriodSeconds: pointy.Int64(0), PropagationPolicy: &policy}

	err := c.Clientset.CoreV1().Namespaces().Delete(context.TODO(), u, opts)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ApplyMetadata(namespace string, fullName string, email string, userType string) error {
	kind := "Namespace"
	apiVersion := "v1"

	// Configure patch
	config := applycorev1.NamespaceApplyConfiguration{
		TypeMetaApplyConfiguration: applymetav1.TypeMetaApplyConfiguration{
			Kind:       &kind,
			APIVersion: &apiVersion,
		},
		ObjectMetaApplyConfiguration: &applymetav1.ObjectMetaApplyConfiguration{
			Name:      &namespace,
			Namespace: &namespace,
			Labels:    map[string]string{LabelUserType: userType},
			Annotations: map[string]string{
				AnnotationUserFullname: fullName,
				AnnotationUserEmail:    email,
			},
		},
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Set field manager to "kubectl" in order
	// to avoid any conflicts or hacks (e.g. having to set Force=true).
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err := c.Clientset.CoreV1().Namespaces().Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "kubectl"},
	)
	if err != nil {
		return err
	}

	return nil
}
