package k8s

import (
	"reflect"
	"testing"

	internalfake "github.com/uitml/quimby/internal/fake"
	"github.com/uitml/quimby/internal/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_resourceAsInt64(t *testing.T) {
	type args struct {
		resources corev1.ResourceList
		names     []corev1.ResourceName
	}
	tests := []struct {
		name    string
		args    args
		want    map[corev1.ResourceName]int64
		wantErr bool
	}{
		// Testcase 1: Multiple input, multiple output + decimal value
		{
			name: "MIMO",
			args: args{
				resources: internalfake.NewResourceQuota(
					"foo123", // namespace
					4500,     // cpu
					2,        // gpu
					16,       // memory
					2,        // inverse scaling
				).Status.Hard,
				names: []corev1.ResourceName{
					ResourceRequestsGPU,
					corev1.ResourceRequestsCPU,
					corev1.ResourceRequestsMemory,
				},
			},
			want: map[corev1.ResourceName]int64{
				ResourceRequestsGPU:           2,
				corev1.ResourceRequestsCPU:    4500,
				corev1.ResourceRequestsMemory: (16*1024 + 256) * 1024 * 1024, // bytes
			},
			wantErr: false,
		},
		// Testcase 2: One or more resources doesn't exist. Should return nil and error
		{
			name: "Nonexistend resource",
			args: args{
				resources: internalfake.NewResourceQuota(
					"foo123", // namespace
					4500,     // cpu
					2,        // gpu
					16,       // memory
					2,        // inverse scaling
				).Status.Hard,
				names: []corev1.ResourceName{
					ResourceRequestsGPU,
					"foo",
					corev1.ResourceRequestsMemory,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resourceAsInt64(tt.args.resources, tt.args.names...)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceAsInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceAsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestClient_Quota(t *testing.T) {
	type fields struct {
		Clientset kubernetes.Interface
	}
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    resource.Quota
		wantErr bool
	}{
		// Testcase 1:
		{
			name: "User exists and has resources",
			fields: fields{
				Clientset: fake.NewSimpleClientset(
					internalfake.NewResourceQuotaList("foo123", 4500, 2, 16, 2),
					internalfake.NewPVCList("foo123", 500),
				),
			},
			args: args{namespace: "foo123"},
			want: resource.Quota{
				CPU:     resource.Summary{Max: 4500000, Used: 2250000},
				GPU:     resource.Summary{Max: 2, Used: 1},
				Memory:  resource.Summary{Max: (16*1024 + 256) * 1024 * 1024, Used: (16*1024 + 256) * 1024 * 512},
				Storage: 500 * 1024 * 1024 * 1024,
			},
			wantErr: false,
		},
		{
			name: "User does not exist or has no resource quota",
			fields: fields{
				Clientset: fake.NewSimpleClientset(
					internalfake.NewResourceQuotaList("foo123", 4500, 2, 16, 2),
					internalfake.NewPVCList("foo123", 500),
				),
			},
			args:    args{namespace: "bar123"},
			want:    resource.Quota{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Clientset: tt.fields.Clientset,
			}
			got, err := c.Quota(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Quota() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Quota() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_TotalGPUs(t *testing.T) {
	type fields struct {
		Clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    resource.Summary
		wantErr bool
	}{
		// Testcase 1: Empty node list. Should return zero
		{
			name:    "No nodes",
			fields:  fields{fake.NewSimpleClientset()},
			want:    resource.Summary{Max: 0, Used: 0},
			wantErr: false,
		},
		// Testcase 2: 3 servers, 2 with GPUs
		{
			name: "3 srv, 2 with GPU",
			fields: fields{fake.NewSimpleClientset(
				internalfake.NewNodeList([]string{"foo", "bar", "baz"}, []int64{0, 8, 7}, []bool{false, false, false}),
			)},
			want:    resource.Summary{Max: 15, Used: 0},
			wantErr: false,
		},
		// Testcase 3: 3 servers, 2 with GPUs, but one is unschedulable
		{
			name: "3 srv, 2 with gpu, one unschedulable",
			fields: fields{fake.NewSimpleClientset(
				internalfake.NewNodeList([]string{"foo", "bar", "baz"}, []int64{0, 8, 7}, []bool{false, false, true}),
			)},
			want:    resource.Summary{Max: 8, Used: 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Clientset: tt.fields.Clientset,
			}
			got, err := c.TotalGPUs()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TotalGPUs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.TotalGPUs() = %v, want %v", got, tt.want)
			}
		})
	}
}
