package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	applyappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applyrbacv1 "k8s.io/client-go/applyconfigurations/rbac/v1"
)

func (c *Client) Apply(namespace string, manifest []byte) error {
	dec := k8syaml.NewYAMLToJSONDecoder(bytes.NewReader(manifest))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}

		res = append(res, valueBytes)
	}

	// Configure patches
	for _, m := range res {
		if strings.Contains(string(m), "kind: Namespace") {
			err := c.applyNamespace(namespace, m)
			if err != nil {
				return err
			}
		} else if strings.Contains(string(m), "kind: RoleBinding") {
			err := c.applyRoleBinding(namespace, m)
			if err != nil {
				return err
			}
		} else if strings.Contains(string(m), "kind: ResourceQuota") {
			err := c.applyResourceQuota(namespace, m)
			if err != nil {
				return err
			}
		} else if strings.Contains(string(m), "kind: LimitRange") {
			err := c.applyLimitRange(namespace, m)
			if err != nil {
				return err
			}
		} else if strings.Contains(string(m), "kind: PersistentVolumeClaim") {
			err := c.applyPersistentVolumeClaim(namespace, m)
			if err != nil {
				return err
			}
		} else if strings.Contains(string(m), "kind: Deployment") {
			err := c.applyDeployment(namespace, m)
			if err != nil {
				return err
			}
		} else if strings.Contains(string(m), "kind: Service") {
			err := c.applyService(namespace, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) applyNamespace(namespace string, manifest []byte) error {
	config := applycorev1.NamespaceApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.CoreV1().Namespaces().Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}

func (c *Client) applyRoleBinding(namespace string, manifest []byte) error {
	config := applyrbacv1.RoleBindingApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.RbacV1().RoleBindings(namespace).Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}

func (c *Client) applyResourceQuota(namespace string, manifest []byte) error {
	config := applycorev1.ResourceQuotaApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.CoreV1().ResourceQuotas(namespace).Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}

func (c *Client) applyLimitRange(namespace string, manifest []byte) error {
	config := applycorev1.LimitRangeApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.CoreV1().LimitRanges(namespace).Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}

func (c *Client) applyPersistentVolumeClaim(namespace string, manifest []byte) error {
	config := applycorev1.PersistentVolumeClaimApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.CoreV1().PersistentVolumeClaims(namespace).Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}

func (c *Client) applyDeployment(namespace string, manifest []byte) error {
	config := applyappsv1.DeploymentApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.AppsV1().Deployments(namespace).Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}

func (c *Client) applyService(namespace string, manifest []byte) error {
	config := applycorev1.ServiceApplyConfiguration{}
	err := json.Unmarshal(manifest, &config)
	if err != nil {
		return err
	}

	// Using server-side apply.
	// Fields are managed and have an owner. Some kubectl fields are managed
	// by "kubectl", while others are managed by "kubectl-client-side-apply"...
	// Set Force=true to circumvent that.
	// See https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
	_, err = c.Clientset.CoreV1().Services(namespace).Apply(
		context.TODO(),
		&config,
		metav1.ApplyOptions{FieldManager: "quimby", Force: true},
	)

	return err
}
