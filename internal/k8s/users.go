package k8s

import (
	"context"
	"errors"

	"github.com/uitml/quimby/internal/validate"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
