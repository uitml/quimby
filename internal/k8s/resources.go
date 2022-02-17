package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/openlyinc/pointy"
	"github.com/uitml/quimby/internal/resource"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceRequestsGPU corev1.ResourceName = "requests.nvidia.com/gpu"
	ResourceGPU         corev1.ResourceName = "nvidia.com/gpu"
)

func resourceAsInt64(resources corev1.ResourceList, names ...corev1.ResourceName) (map[corev1.ResourceName]int64, error) {
	result := make(map[corev1.ResourceName]int64)

	for _, name := range names {
		// Should probably check if the resource exists
		res, ok := resources[name]

		if !ok {
			return nil, fmt.Errorf("in resourceAsInt64: Resource %v does not exist", name)
		}
		if strings.Contains(string(name), "cpu") {
			// use millivalue
			result[name] = res.MilliValue()
			continue
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

func (c *Client) Quota(namespace string) (resource.Quota, error) {
	// Compute
	res, err := c.Clientset.CoreV1().ResourceQuotas(namespace).Get(context.TODO(), "compute-resources", metav1.GetOptions{})
	if err != nil {
		return resource.Quota{}, err
	}

	// Storage
	pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), "storage", metav1.GetOptions{})
	if err != nil {
		return resource.Quota{}, err
	}

	// Convert all resources to Int64
	maxResources, err := resourceAsInt64(
		res.Spec.Hard,
		ResourceRequestsGPU,
		corev1.ResourceRequestsCPU,
		corev1.ResourceRequestsMemory,
	)
	if err != nil {
		return resource.Quota{}, err
	}

	usedResources, err := resourceAsInt64(
		res.Status.Used,
		ResourceRequestsGPU,
		corev1.ResourceRequestsCPU,
		corev1.ResourceRequestsMemory,
	)
	if err != nil {
		return resource.Quota{}, err
	}

	storage, err := resourceAsInt64(
		pvc.Spec.Resources.Requests,
		corev1.ResourceStorage,
	)
	if err != nil {
		return resource.Quota{}, err
	}

	rq := resource.Quota{
		GPU: resource.Summary{
			Max:  maxResources[ResourceRequestsGPU],
			Used: usedResources[ResourceRequestsGPU],
		},
		CPU: resource.Summary{
			Max:  maxResources[corev1.ResourceRequestsCPU],
			Used: usedResources[corev1.ResourceRequestsCPU],
		},
		Memory: resource.Summary{
			Max:  maxResources[corev1.ResourceRequestsMemory],
			Used: usedResources[corev1.ResourceRequestsMemory],
		},
		Storage: storage[corev1.ResourceStorage],
	}

	return rq, nil
}

func (c *Client) Spec(namespace string) (*resource.Spec, error) {
	errchan := make(chan error)

	// Compute
	reschan := make(chan *corev1.ResourceQuota)
	go func() {
		res, err := c.Clientset.CoreV1().ResourceQuotas(namespace).Get(context.TODO(), "compute-resources", metav1.GetOptions{})
		if err != nil {
			errchan <- err
		}
		reschan <- res
	}()

	// Limits
	limchan := make(chan *corev1.LimitRange)
	go func() {
		lim, err := c.Clientset.CoreV1().LimitRanges(namespace).Get(context.TODO(), "default-resources", metav1.GetOptions{})
		if err != nil {
			errchan <- err
		}
		limchan <- lim
	}()

	// Storage
	pvcchan := make(chan *corev1.PersistentVolumeClaim)
	go func() {
		pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), "storage", metav1.GetOptions{})
		if err != nil {
			errchan <- err
		}
		pvcchan <- pvc
	}()

	// storage-proxy
	proxychan := make(chan *appsv1.Deployment)
	go func() {
		dpl, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.TODO(), "storage-proxy", metav1.GetOptions{})
		if err != nil {
			errchan <- err
		}
		proxychan <- dpl
	}()

	var maxResources, defaultLimits, storage, proxy, proxyrequest map[corev1.ResourceName]int64
	var err error
	// Receive and convert all resources to Int64
	for i := 0; i < 4; i++ {
		select {
		case res := <-reschan:
			maxResources, err = resourceAsInt64(
				res.Spec.Hard,
				ResourceRequestsGPU,
				corev1.ResourceRequestsCPU,
				corev1.ResourceRequestsMemory,
			)
			if err != nil {
				return nil, err
			}

			continue
		case lim := <-limchan:
			defaultLimits, err = resourceAsInt64(
				lim.Spec.Limits[0].Default,
				corev1.ResourceCPU,
				corev1.ResourceMemory,
				ResourceGPU,
			)
			if err != nil {
				return nil, err
			}

			continue
		case pvc := <-pvcchan:
			storage, err = resourceAsInt64(
				pvc.Spec.Resources.Requests,
				corev1.ResourceStorage,
			)
			if err != nil {
				return nil, err
			}

			continue
		case dpl := <-proxychan:
			proxy, err = resourceAsInt64(
				dpl.Spec.Template.Spec.Containers[0].Resources.Limits,
				corev1.ResourceCPU,
				corev1.ResourceMemory,
			)
			if err != nil {
				return nil, err
			}

			proxyrequest, err = resourceAsInt64(
				dpl.Spec.Template.Spec.Containers[0].Resources.Requests,
				corev1.ResourceCPU,
			)
			if err != nil {
				return nil, err
			}

			continue
		case err := <-errchan:
			close(errchan)
			return nil, err
		}
	}

	// This is a hot mess
	result := resource.Spec{
		GPU:                    pointy.Int64(maxResources[ResourceRequestsGPU]),
		GPUPerJob:              pointy.Int64(defaultLimits[ResourceGPU]),
		MaxMemoryPerJob:        pointy.Int64(maxResources[corev1.ResourceRequestsMemory] / 1024 / 1024 / 1024 / maxResources[ResourceRequestsGPU]), // in GiB
		DefaultMemoryPerJob:    pointy.Int64(defaultLimits[corev1.ResourceMemory] / 1024 / 1024 / 1024),                                            // in GiB
		CPUPerJob:              pointy.Int64(defaultLimits[corev1.ResourceCPU] / 1000),                                                             // not milli
		StorageProxyCPURequest: pointy.Int64(proxyrequest[corev1.ResourceCPU]),                                                                     // milli
		StorageProxyCPULimit:   pointy.Int64(proxy[corev1.ResourceCPU]),                                                                            // milli
		StorageProxyMemory:     pointy.Int64(proxy[corev1.ResourceMemory] / 1024 / 1024),                                                           // in MB
		StorageSize:            pointy.Int64(storage[corev1.ResourceStorage] / 1024 / 1024 / 1024),                                                 // in GiB
	}

	return &result, nil
}

func (c *Client) DefaultRequest(namespace string) (resource.Request, error) {
	res, err := c.Clientset.CoreV1().LimitRanges(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return resource.Request{}, err
	}

	limits, err := resourceAsInt64(
		res.Items[0].Spec.Limits[0].DefaultRequest,
		ResourceGPU,
		corev1.ResourceCPU,
		corev1.ResourceRequestsMemory,
	)
	if err != nil {
		return resource.Request{}, err
	}

	rr := resource.Request{
		GPU:    limits[ResourceGPU],
		CPU:    limits[corev1.ResourceCPU],
		Memory: limits[corev1.ResourceMemory],
	}

	return rr, nil
}

func (c *Client) TotalGPUs() (resource.Summary, error) {
	var totalGPUs int64 = 0
	// TODO: Find out how to get used GPUs from node info. Haven't found it yet, so now it's being counted from the user info.
	var usedGPUs int64 = 0

	nodes, err := c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return resource.Summary{}, err
	}

	for _, node := range nodes.Items {
		var nodeReady bool = true
		for _, condition := range node.Status.Conditions {
			if condition.Type != corev1.NodeReady {
				continue
			}
			if condition.Status != corev1.ConditionTrue {
				nodeReady = false
			}
			break
		}
		if !nodeReady {
			continue
		}
		if node.Spec.Unschedulable {
			continue
		}

		// Ignoring errors here since some nodes might not have all resources
		g, _ := resourceAsInt64(node.Status.Capacity, ResourceGPU)
		totalGPUs += g[ResourceGPU]
	}

	return resource.Summary{Max: totalGPUs, Used: usedGPUs}, nil
}
