package user

func memoryPerGPU(usr User) int64 {
	memory := usr.ResourceQuota.Memory.Max
	GPUMax := usr.ResourceQuota.GPU.Max
	memoryPerGPU := memory / GPUMax

	return memoryPerGPU
}

func TotalResourcesUsed(userList []User) map[string]int {
	var r = make(map[string]int)

	r["GPU"] = 0

	for _, usr := range userList {
		r["GPU"] += int(usr.ResourceQuota.GPU.Used)
	}

	return r
}
