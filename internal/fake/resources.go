package fake

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceRequestsGPU corev1.ResourceName = "requests.nvidia.com/gpu"
	ResourceGPU         corev1.ResourceName = "nvidia.com/gpu"
)

func NewResourceQuotaList(namespace string, cpu int64, gpu int64, memory int64, inverseScaling int64) *corev1.ResourceQuotaList {
	quota := corev1.ResourceQuotaList{
		TypeMeta: metav1.TypeMeta{Kind: "ResourceQuotaList", APIVersion: "v1"},
		Items:    []corev1.ResourceQuota{*NewResourceQuota(namespace, cpu, gpu, memory, inverseScaling)},
	}

	return &quota
}

func NewResourceQuota(namespace string, cpu int64, gpu int64, memory int64, inverseScaling int64) *corev1.ResourceQuota {
	quota := corev1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{Kind: "ResourceQuota", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "compute-resources",
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceRequestsCPU:    *resource.NewQuantity(cpu, resource.DecimalSI),
				ResourceRequestsGPU:           *resource.NewQuantity(gpu, resource.DecimalSI),
				corev1.ResourceRequestsMemory: *resource.NewQuantity((memory*1024+256)*1024*1024, resource.BinarySI),
			},
		},
		Status: corev1.ResourceQuotaStatus{
			Hard: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceRequestsCPU:    *resource.NewMilliQuantity(cpu, resource.DecimalSI),
				ResourceRequestsGPU:           *resource.NewQuantity(gpu, resource.DecimalSI),
				corev1.ResourceRequestsMemory: *resource.NewQuantity((memory*1024+256)*1024*1024, resource.BinarySI),
			},
			Used: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceRequestsCPU:    *resource.NewQuantity(cpu/inverseScaling, resource.DecimalSI),
				ResourceRequestsGPU:           *resource.NewQuantity(gpu/inverseScaling, resource.DecimalSI),
				corev1.ResourceRequestsMemory: *resource.NewQuantity((memory*1024+256)*1024*1024/inverseScaling, resource.BinarySI),
			},
		},
	}

	return &quota
}

func NewPVCList(namespace string, size int64) *corev1.PersistentVolumeClaimList {
	quota := corev1.PersistentVolumeClaimList{
		TypeMeta: metav1.TypeMeta{Kind: "PersistentVolumeClaimList", APIVersion: "v1"},
		Items:    []corev1.PersistentVolumeClaim{*NewPVC(namespace, size)},
	}

	return &quota
}

func NewPVC(namespace string, size int64) *corev1.PersistentVolumeClaim {
	storageClass := "nfs-storage"
	quota := corev1.PersistentVolumeClaim{
		TypeMeta:   metav1.TypeMeta{Kind: "PersistentVolumeClaim", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "storage", Namespace: namespace},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteMany"},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: *resource.NewQuantity(
						size*1024*1024*1024,
						resource.BinarySI,
					),
				},
			},
			VolumeName:       "storage",
			StorageClassName: &storageClass,
		},
	}

	return &quota
}

func NewNodeList(nodeNames []string, gpus []int64, unschedulable []bool) *corev1.NodeList {
	var nodes []corev1.Node
	for i, name := range nodeNames {
		nodes = append(nodes, newNode(name, gpus[i], unschedulable[i]))
	}

	nodeList := corev1.NodeList{
		TypeMeta: metav1.TypeMeta{Kind: "NodeList", APIVersion: "v1"},
		Items:    nodes,
	}

	return &nodeList
}

func newNode(name string, gpus int64, isUnschedulable bool) corev1.Node {
	capacity := map[corev1.ResourceName]resource.Quantity{}
	if gpus > 0 {
		capacity = map[corev1.ResourceName]resource.Quantity{ResourceGPU: *resource.NewQuantity(gpus, resource.DecimalSI)}
	}

	node := corev1.Node{
		TypeMeta:   metav1.TypeMeta{Kind: "Node", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.NodeSpec{Unschedulable: isUnschedulable},
		Status:     corev1.NodeStatus{Capacity: capacity, Allocatable: capacity},
	}

	return node
}
