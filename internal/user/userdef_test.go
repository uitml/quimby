/*
This package implements tools and data structures for operating on users.
*/

package user

import (
	"reflect"
	"testing"

	"github.com/uitml/quimby/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		// The only behaviour that should if fields are missing is that e-mail is auto-generated.
		{
			name: "all fields present",
			args: args{
				namespace: corev1.Namespace{
					metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
					metav1.ObjectMeta{
						Name:        "fba000",
						Labels:      map[string]string{k8s.LabelUserType: "admin"},
						Annotations: map[string]string{k8s.AnnotationUserFullname: "Foo Bar", k8s.AnnotationUserEmail: "foo@bar.baz"},
					},
					corev1.NamespaceSpec{},
					corev1.NamespaceStatus{},
				},
			},
			want: User{Username: "fba000", Fullname: "Foo Bar", Email: "foo@bar.baz", Usertype: "admin"},
		},
		{
			name: "missing annotations and labels",
			args: args{
				namespace: corev1.Namespace{
					metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
					metav1.ObjectMeta{
						Name:        "boo001",
						Labels:      map[string]string{},
						Annotations: map[string]string{},
					},
					corev1.NamespaceSpec{},
					corev1.NamespaceStatus{},
				},
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
