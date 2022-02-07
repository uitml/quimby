/*
This package implements tools and data structures for operating on users.
*/

package user

import (
	"reflect"
	"testing"

	"github.com/uitml/quimby/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func newFakeNamespace(name string, labels map[string]string, annotations map[string]string) *corev1.Namespace {
	ns := corev1.Namespace{
		metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		corev1.NamespaceSpec{},
		corev1.NamespaceStatus{},
	}

	return &ns
}

func newFakeResourceQuota(namespace string, cpu int64, gpu int64, memory int64, inverseScaling int64) *corev1.ResourceQuota {
	quota := corev1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{Kind: "ResourceQuota", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "compute-resources",
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceRequestsCPU:    *resource.NewQuantity(cpu, resource.DecimalSI),
				k8s.ResourceRequestsGPU:       *resource.NewQuantity(gpu, resource.DecimalSI),
				corev1.ResourceRequestsMemory: *resource.NewQuantity((memory*1024+256)*1024*1024, resource.BinarySI),
			},
		},
		Status: corev1.ResourceQuotaStatus{
			Hard: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceRequestsCPU:    *resource.NewMilliQuantity(cpu, resource.DecimalSI),
				k8s.ResourceRequestsGPU:       *resource.NewQuantity(gpu, resource.DecimalSI),
				corev1.ResourceRequestsMemory: *resource.NewQuantity((memory*1024+256)*1024*1024, resource.BinarySI),
			},
			Used: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceRequestsCPU:    *resource.NewQuantity(cpu/inverseScaling, resource.DecimalSI),
				k8s.ResourceRequestsGPU:       *resource.NewQuantity(gpu/inverseScaling, resource.DecimalSI),
				corev1.ResourceRequestsMemory: *resource.NewQuantity((memory*1024+256)*1024*1024/inverseScaling, resource.BinarySI),
			},
		},
	}

	return &quota
}

func newFakeResourceQuotaList(namespace string, cpu int64, gpu int64, memory int64, inverseScaling int64) *corev1.ResourceQuotaList {
	quota := corev1.ResourceQuotaList{
		TypeMeta: metav1.TypeMeta{Kind: "ResourceQuotaList", APIVersion: "v1"},
		Items:    []corev1.ResourceQuota{*newFakeResourceQuota(namespace, cpu, gpu, memory, inverseScaling)},
	}

	return &quota
}

func newFakePVC(namespace string, size int64, inverseScaling int64) *corev1.PersistentVolumeClaim {
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

func newFakePVCList(namespace string, size int64, inverseScaling int64) *corev1.PersistentVolumeClaimList {
	quota := corev1.PersistentVolumeClaimList{
		TypeMeta: metav1.TypeMeta{Kind: "PersistentVolumeClaimList", APIVersion: "v1"},
		Items:    []corev1.PersistentVolumeClaim{*newFakePVC(namespace, size, inverseScaling)},
	}

	return &quota
}

func TestFromNamespace(t *testing.T) {
	type args struct {
		namespace corev1.Namespace
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		// The only behaviour that should change if fields are missing is that e-mail is auto-generated.
		{
			name: "all fields present",
			args: args{
				namespace: *newFakeNamespace(
					"fba000",
					map[string]string{k8s.LabelUserType: "admin"},
					map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
				),
			},
			want: User{Username: "fba000", Fullname: "Foo Bar", Email: "foo@bar.baz", Usertype: "admin"},
		},
		{
			name: "missing annotations and labels",
			args: args{
				namespace: *newFakeNamespace(
					"boo001",
					map[string]string{},
					map[string]string{},
				),
			},
			want: User{Username: "boo001", Fullname: "", Email: "boo001@post.uit.no", Usertype: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromNamespace(tt.args.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPopulateList(t *testing.T) {
	type args struct {
		c             k8s.ResourceClient
		listResources bool
	}
	tests := []struct {
		name    string
		args    args
		want    []User
		wantErr bool
	}{
		// TODO: Add test cases.
		// 1. Normal usecase with listed resources
		// 2. Normal usecase without
		// 3. No users in cluster
		{
			name: "all fields present - list resources",
			args: args{
				c: &k8s.Client{Clientset: fake.NewSimpleClientset(
					newFakeNamespace(
						"foo123",
						map[string]string{k8s.LabelUserType: "admin"},
						map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
					),
					newFakeResourceQuotaList("foo123", 4500, 2, 16, 2),
					newFakePVCList("foo123", 500, 2),
				)},
				listResources: true,
			},
			want: []User{
				{
					Username: "foo123",
					Fullname: "Foo Bar",
					Email:    "foo@bar.baz",
					Usertype: "admin",
					ResourceQuota: k8s.ResourceQuota{
						CPU:     k8s.ResourceSummary{Max: 4500, Used: 2250},
						GPU:     k8s.ResourceSummary{Max: 2, Used: 1},
						Memory:  k8s.ResourceSummary{Max: (16*1024 + 256) * 1024 * 1024, Used: (16*1024 + 256) * 1024 * 512},
						Storage: 500 * 1024 * 1024 * 1024,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "all fields present - no resources",
			args: args{
				c: &k8s.Client{Clientset: fake.NewSimpleClientset(
					newFakeNamespace(
						"foo123",
						map[string]string{k8s.LabelUserType: "admin"},
						map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
					),
					newFakeResourceQuotaList("foo123", 4500, 2, 16, 2),
					newFakePVCList("foo123", 500, 2),
				)},
				listResources: false,
			},
			want: []User{
				{
					Username:      "foo123",
					Fullname:      "Foo Bar",
					Email:         "foo@bar.baz",
					Usertype:      "admin",
					ResourceQuota: k8s.ResourceQuota{},
				},
			},
			wantErr: false,
		},
		{
			name: "No users in cluster",
			args: args{
				c:             &k8s.Client{Clientset: fake.NewSimpleClientset()},
				listResources: false,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PopulateList(tt.args.c, tt.args.listResources)
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulateList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PopulateList() = %v, want %v", got, tt.want)
			}
		})
	}
}
