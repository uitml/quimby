package k8s

import (
	"context"
	"errors"
	"fmt"

	"github.com/uitml/quimby/internal/validate"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const (
	AnnotationUserFullname string = "springfield.uit.no/user-fullname"
	AnnotationUserEmail    string = "springfield.uit.no/user-email"
	LabelUserType          string = "springfield.uit.no/user-type"
)

func (c *Client) GetNamespaceList() (*corev1.NamespaceList, error) {
	namespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaceList, nil
}

func (c *Client) Namespace(username string) (*corev1.Namespace, error) {
	namespaces, err := c.GetNamespaceList()
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

func (c *Client) NewUser(username string, fullname string, email string, userType string) error {
	if !validate.Username(username) {
		return errors.New("invalid username. Username must be on the form xyz123")
	}

	labels := map[string]string{LabelUserType: validate.DefaultIfEmpty(userType, "student")}
	annotations := map[string]string{AnnotationUserFullname: fullname, AnnotationUserEmail: email}

	ns := NewNamespace(username, labels, annotations)

	// Create namespace
	_, err := c.Clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})

	// Resources

	// Resource Limits

	// Storage

	return err
}

func (c *Client) NewSimpleUser(username string) error {
	return c.NewUser(username, "", "", "")
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

func newRoleBinding(username string) *rbacv1.RoleBinding {
	rb := rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "owner",
			Namespace: username,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     "User",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "springfield:namespace-owner",
		},
	}

	return &rb
}

func newLimitRange(username string, defaultCPU string, defaultMemory string, defaultGPU int64) (*corev1.LimitRange, error) {
	rCPU, err := resource.ParseQuantity(defaultCPU)
	if err != nil {
		return nil, err
	}

	rMemory, err := resource.ParseQuantity(defaultMemory)
	if err != nil {
		return nil, err
	}

	lr := corev1.LimitRange{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LimitRange",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-resources",
			Namespace: username,
		},
		Spec: corev1.LimitRangeSpec{
			Limits: []corev1.LimitRangeItem{
				{
					Type: corev1.LimitTypeContainer,
					Default: corev1.ResourceList{
						corev1.ResourceCPU:    rCPU,
						corev1.ResourceMemory: rMemory,
						ResourceGPU:           *resource.NewQuantity(defaultGPU, resource.DecimalSI),
					},
				},
			},
		},
	}

	return &lr, nil
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

func (c *Client) ApplyMetaData(namespace string, fullName string, email string, userType string) error {
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
