package k8s

import (
	"reflect"
	"testing"

	internalfake "github.com/uitml/quimby/internal/fake"
	corev1 "k8s.io/api/core/v1"
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
				ResourceRequestsGPU: 2,
				// Fake resource quota sets CPU to a milli value (decimal) on Status.Hard, so it's rounded up to 5.
				corev1.ResourceRequestsCPU:    5,
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
