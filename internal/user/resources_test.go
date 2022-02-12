package user

import (
	"reflect"
	"testing"

	"github.com/uitml/quimby/internal/k8s"
	corev1 "k8s.io/api/core/v1"
)

func newFakeResourceUser(CPUMax int64, CPUUsed int64, memoryMax int64, memoryUsed int64, GPUMax int64, GPUUsed int64, storage int64) User {
	usr := User{
		Username: "",
		fullname: "",
		email:    "",
		usertype: "",
		ResourceQuota: k8s.ResourceQuota{
			CPU:     k8s.ResourceSummary{Max: CPUMax, Used: CPUUsed},
			Memory:  k8s.ResourceSummary{Max: memoryMax * 1024 * 1024 * 1024, Used: memoryUsed},
			GPU:     k8s.ResourceSummary{Max: GPUMax, Used: GPUUsed},
			Storage: storage,
		},
	}

	return usr
}

func Test_memoryPerGPU(t *testing.T) {
	type args struct {
		usr User
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// Normal usecase
		{
			name: "Normal usecase",
			args: args{
				usr: newFakeResourceUser(
					0,  // Max CPU
					0,  // Used CPU
					32, // Max memory
					0,  // Used memory
					2,  // Max GPUs
					0,  // Used GPUs
					0,  // Storage
				),
			},
			want: 16 * 1024 * 1024 * 1024,
		},
		// No GPUs assigned
		{
			name: "Normal usecase",
			args: args{
				usr: newFakeResourceUser(
					0,  // Max CPU
					0,  // Used CPU
					32, // Max memory
					0,  // Used memory
					0,  // Max GPUs
					0,  // Used GPUs
					0,  // Storage
				),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := memoryPerGPU(tt.args.usr); got != tt.want {
				t.Errorf("memoryPerGPU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTotalResourcesUsed(t *testing.T) {
	type args struct {
		userList []User
	}
	tests := []struct {
		name string
		args args
		want map[corev1.ResourceName]int64
	}{
		// Test case 1: Empty user list (default to zeros for all vals)
		{
			name: "Empty user list",
			args: args{
				userList: []User{},
			},
			want: map[corev1.ResourceName]int64{
				k8s.ResourceGPU:        0,
				corev1.ResourceCPU:     0,
				corev1.ResourceMemory:  0,
				corev1.ResourceStorage: 0,
			},
		},
		// Test case 2: Empty resource for all users (default to zeros for all vals)
		{
			name: "Empty resources",
			args: args{
				userList: []User{
					{
						Username: "foo123",
						fullname: "Foo Bar",
						email:    "foo@bar.baz",
						usertype: "alumni",
					},
					{
						Username: "bar321",
						fullname: "Bar Baz",
						email:    "bar@baz.com",
						usertype: "student",
					},
				},
			},
			want: map[corev1.ResourceName]int64{
				k8s.ResourceGPU:        0,
				corev1.ResourceCPU:     0,
				corev1.ResourceMemory:  0,
				corev1.ResourceStorage: 0,
			},
		},
		// Test case 3: Empty resource for some users
		{
			name: "Semi-empty resources",
			args: args{
				userList: []User{
					newFakeResourceUser(
						4500, // Max CPU
						2500, // Used CPU
						32,   // Max memory
						16,   // Used memory
						2,    // Max GPUs
						1,    // Used GPUs
						500,  // Storage
					),
					{
						Username: "bar321",
						fullname: "Bar Baz",
						email:    "bar@baz.com",
						usertype: "student",
					},
				},
			},
			want: map[corev1.ResourceName]int64{
				k8s.ResourceGPU:        1,
				corev1.ResourceCPU:     2500,
				corev1.ResourceMemory:  16,
				corev1.ResourceStorage: 500,
			},
		},
		// Test case 4: All users have resources
		{
			name: "Full resource population",
			args: args{
				userList: []User{
					newFakeResourceUser(
						4500, // Max CPU
						2500, // Used CPU
						32,   // Max memory
						16,   // Used memory
						2,    // Max GPUs
						1,    // Used GPUs
						500,  // Storage
					),
					newFakeResourceUser(
						8500, // Max CPU
						4500, // Used CPU
						64,   // Max memory
						32,   // Used memory
						4,    // Max GPUs
						2,    // Used GPUs
						500,  // Storage
					),
				},
			},
			want: map[corev1.ResourceName]int64{
				k8s.ResourceGPU:        3,
				corev1.ResourceCPU:     7000,
				corev1.ResourceMemory:  48,
				corev1.ResourceStorage: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TotalResourcesUsed(tt.args.userList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TotalResourcesUsed() = %v, want %v", got, tt.want)
			}
		})
	}
}
