package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceRequestsGPU corev1.ResourceName = "requests.nvidia.com/gpu"
	ResourceGPU         corev1.ResourceName = "nvidia.com/gpu"
)

type ResourceSummary struct {
	Max  int64
	Used int64
}

type ResourceQuota struct {
	GPU     ResourceSummary
	CPU     ResourceSummary
	Memory  ResourceSummary
	Storage int64
}

type ResourceRequest struct {
	GPU    int64
	CPU    int64
	Memory int64
}

func resourceAsInt64(resources corev1.ResourceList, names ...corev1.ResourceName) (map[corev1.ResourceName]int64, error) {
	result := make(map[corev1.ResourceName]int64)

	for _, name := range names {
		// Should probably check if the resource exists
		res, ok := resources[name]

		if !ok {
			return nil, fmt.Errorf("error in resourceAsInt64: Resource %v does not exist", name)
		}
		val, ok := res.AsInt64()

		if ok {
			result[name] = val
		} else {
			result[name] = res.ToDec().Value()
		}
	}
	return result, nil
}

func (c *Client) GetResourceQuota(namespace string) (ResourceQuota, error) {
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

	// Convert all resources to Int64
	maxResources, err := resourceAsInt64(
		res.Items[0].Spec.Hard,
		ResourceRequestsGPU,
		corev1.ResourceRequestsCPU,
		corev1.ResourceRequestsMemory,
	)
	if err != nil {
		return ResourceQuota{}, err
	}

	usedResources, err := resourceAsInt64(
		res.Items[0].Status.Used,
		ResourceRequestsGPU,
		corev1.ResourceRequestsCPU,
		corev1.ResourceRequestsMemory,
	)
	if err != nil {
		return ResourceQuota{}, err
	}

	storage, err := resourceAsInt64(
		pvc.Items[0].Spec.Resources.Requests,
		corev1.ResourceStorage,
	)
	if err != nil {
		return ResourceQuota{}, err
	}

	rq := ResourceQuota{
		GPU: ResourceSummary{
			Max:  maxResources[ResourceRequestsGPU],
			Used: usedResources[ResourceRequestsGPU],
		},
		CPU: ResourceSummary{
			Max:  maxResources[corev1.ResourceRequestsCPU],
			Used: usedResources[corev1.ResourceRequestsCPU],
		},
		Memory: ResourceSummary{
			Max:  maxResources[corev1.ResourceRequestsMemory],
			Used: usedResources[corev1.ResourceRequestsMemory],
		},
		Storage: storage[corev1.ResourceStorage],
	}

	return rq, nil
}

func (c *Client) GetDefaultRequest(namespace string) (ResourceRequest, error) {
	res, err := c.Clientset.CoreV1().LimitRanges(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ResourceRequest{}, err
	}

	limits, err := resourceAsInt64(
		res.Items[0].Spec.Limits[0].DefaultRequest,
		ResourceGPU,
		corev1.ResourceCPU,
		corev1.ResourceRequestsMemory,
	)
	if err != nil {
		return ResourceRequest{}, err
	}

	rr := ResourceRequest{
		GPU:    limits[ResourceGPU],
		CPU:    limits[corev1.ResourceCPU],
		Memory: limits[corev1.ResourceMemory],
	}

	return rr, nil
}

func (c *Client) GetTotalGPUs() (int64, error) {
	var totalGPUs int64 = 0

	nodes, err := c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	for _, node := range nodes.Items {
		if !node.Spec.Unschedulable {
			// Ignoring errors here since some nodes might not have all resources
			g, _ := resourceAsInt64(node.Status.Capacity, ResourceGPU)

			totalGPUs += g[ResourceGPU]
		}
	}

	return totalGPUs, nil
}
