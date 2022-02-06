package user

import (
	"fmt"
	"strconv"
	"strings"
)

func memoryPerGPU(usr *User) (string, error) {
	memory, err := strconv.Atoi(strings.Trim(usr.ResourceQuota.Memory.Max, "Mi"))

	if err != nil {
		return "", err
	}

	GPUMax, err := strconv.Atoi(usr.ResourceQuota.GPU.Max)

	if err != nil {
		return "", err
	}

	memoryPerGPU := memory / (GPUMax * 1024)

	return fmt.Sprint(memoryPerGPU) + "Gi", nil
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
