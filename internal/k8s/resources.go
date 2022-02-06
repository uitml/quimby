package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type resourceSummary struct {
	Max  string
	Used string
}

type ResourceQuota struct {
	CPU     resourceSummary
	GPU     resourceSummary
	Memory  resourceSummary
	Storage string
}

type ResourceRequest struct {
	CPU    string
	GPU    string
	Memory string
}

func (c *K8sClient) GetResourceQuota(namespace string) (ResourceQuota, error) {
	// Compute
	res, err := c.Clientset.CoreV1().ResourceQuotas(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ResourceQuota{}, err
	}

	// Storage
	pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ResourceQuota{}, err
	}

	rq := ResourceQuota{
		CPU: resourceSummary{
			Max:  fmt.Sprint(res.Items[0].Spec.Hard["requests.cpu"].ToUnstructured()),
			Used: fmt.Sprint(res.Items[0].Status.Used["requests.cpu"].ToUnstructured()),
		},
		GPU: resourceSummary{
			Max:  fmt.Sprint(res.Items[0].Spec.Hard["requests.nvidia.com/gpu"].ToUnstructured()),
			Used: fmt.Sprint(res.Items[0].Status.Used["requests.nvidia.com/gpu"].ToUnstructured()),
		},
		Memory: resourceSummary{
			Max:  fmt.Sprint(res.Items[0].Spec.Hard["requests.memory"].ToUnstructured()),
			Used: fmt.Sprint(res.Items[0].Status.Used["requests.memory"].ToUnstructured()),
		},
		Storage: fmt.Sprint(pvc.Items[0].Spec.Resources.Requests["storage"].ToUnstructured()),
	}

	return rq, nil
}

func (c *K8sClient) GetDefaultRequest(namespace string) (ResourceRequest, error) {
	res, err := c.Clientset.CoreV1().LimitRanges(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ResourceRequest{}, err
	}

	rr := ResourceRequest{
		CPU:    fmt.Sprint(res.Items[0].Spec.Limits[0].DefaultRequest["cpu"].ToUnstructured()),
		GPU:    fmt.Sprint(res.Items[0].Spec.Limits[0].DefaultRequest["nvidia.com/gpu"].ToUnstructured()),
		Memory: fmt.Sprint(res.Items[0].Spec.Limits[0].DefaultRequest["memory"].ToUnstructured()),
	}

	return rr, nil
}
