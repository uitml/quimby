package user

import (
	"github.com/uitml/quimby/internal/k8s"
	corev1 "k8s.io/api/core/v1"
)

func memoryPerGPU(usr User) int64 {
	// I don't want the program to panic just because some resources on the cluster
	// is set to zero for one user.
	if usr.ResourceQuota.GPU.Max == 0 {
		return 0
	}
	return usr.ResourceQuota.Memory.Max / usr.ResourceQuota.GPU.Max
}

func TotalResourcesUsed(userList []User) map[corev1.ResourceName]int64 {
	r := map[corev1.ResourceName]int64{k8s.ResourceGPU: 0, corev1.ResourceCPU: 0, corev1.ResourceMemory: 0, corev1.ResourceStorage: 0}

	for _, usr := range userList {
		r[k8s.ResourceGPU] += usr.ResourceQuota.GPU.Used
		r[corev1.ResourceCPU] += usr.ResourceQuota.CPU.Used
		r[corev1.ResourceMemory] += usr.ResourceQuota.Memory.Used
		r[corev1.ResourceStorage] += usr.ResourceQuota.Storage
	}

	return r
}
