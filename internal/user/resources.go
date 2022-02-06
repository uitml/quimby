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

func TotalResourcesUsed(userList []User) (map[string]int, error) {
	var r = make(map[string]int)
	var tempGPU int = 0
	var err error

	r["GPU"] = 0

	for _, usr := range userList {
		tempGPU, err = strconv.Atoi(usr.ResourceQuota.GPU.Used)
		if err != nil {
			return nil, err
		}
		r["GPU"] += tempGPU
	}

	return r, nil
}
