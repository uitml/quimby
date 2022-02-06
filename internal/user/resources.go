package user

import (
	"fmt"
	"strconv"
	"strings"
)

func memoryPerGPU(usr *User) string {
	memory, err := strconv.Atoi(strings.Trim(usr.ResourceQuota.Memory.Max, "Mi"))

	if err != nil {
		panic(err)
	}

	GPUMax, err := strconv.Atoi(usr.ResourceQuota.GPU.Max)

	if err != nil {
		panic(err)
	}

	memoryPerGPU := memory / (GPUMax * 1024)

	return fmt.Sprint(memoryPerGPU) + "Gi"
}
