/*
This package implements tools and data structures for operating on users.
*/

package user

import (
	"reflect"
	"testing"

	internalfake "github.com/uitml/quimby/internal/fake"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

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
				namespace: *k8s.NewNamespace(
					"fba000",
					map[string]string{k8s.LabelUserType: "admin"},
					map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
				),
			},
			want: User{Username: "fba000", fullname: "Foo Bar", email: "foo@bar.baz", usertype: "admin"},
		},
		{
			name: "missing annotations and labels",
			args: args{
				namespace: *k8s.NewNamespace(
					"boo001",
					map[string]string{},
					map[string]string{},
				),
			},
			want: User{Username: "boo001", fullname: "", email: "boo001@post.uit.no", usertype: ""},
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
					k8s.NewNamespace(
						"foo123",
						map[string]string{k8s.LabelUserType: "admin"},
						map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
					),
					internalfake.NewResourceQuotaList("foo123", 4500, 2, 16, 2),
					internalfake.NewPVCList("foo123", 500),
				)},
				listResources: true,
			},
			want: []User{
				{
					Username: "foo123",
					fullname: "Foo Bar",
					email:    "foo@bar.baz",
					usertype: "admin",
					ResourceQuota: resource.Quota{
						CPU:     resource.Summary{Max: 4500000, Used: 2250000},
						GPU:     resource.Summary{Max: 2, Used: 1},
						Memory:  resource.Summary{Max: (16*1024 + 256) * 1024 * 1024, Used: (16*1024 + 256) * 1024 * 512},
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
					k8s.NewNamespace(
						"foo123",
						map[string]string{k8s.LabelUserType: "admin"},
						map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
					),
					internalfake.NewResourceQuotaList("foo123", 4500, 2, 16, 2),
					internalfake.NewPVCList("foo123", 500),
				)},
				listResources: false,
			},
			want: []User{
				{
					Username:      "foo123",
					fullname:      "Foo Bar",
					email:         "foo@bar.baz",
					usertype:      "admin",
					ResourceQuota: resource.Quota{},
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
