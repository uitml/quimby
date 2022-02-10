package k8s

import (
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClient_UserExists(t *testing.T) {
	type fields struct {
		Clientset kubernetes.Interface
	}
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// Testcase 1: User exists
		{
			name: "User exists",
			fields: fields{
				Clientset: fake.NewSimpleClientset(
					NewNamespace(
						"foo123",
						map[string]string{},
						map[string]string{},
					),
				),
			},
			args:    args{u: "foo123"},
			want:    true,
			wantErr: false,
		},
		// Testcase 2: User does not exist
		{
			name: "User does not exist",
			fields: fields{
				Clientset: fake.NewSimpleClientset(
					NewNamespace(
						"foo123",
						map[string]string{},
						map[string]string{},
					),
				),
			},
			args:    args{u: "bar113"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Clientset: tt.fields.Clientset,
			}
			got, err := c.UserExists(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.UserExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	type fields struct {
		Clientset kubernetes.Interface
	}
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Testcase 1: User exists
		{
			name: "User exists",
			fields: fields{
				Clientset: fake.NewSimpleClientset(
					NewNamespace(
						"foo123",
						map[string]string{},
						map[string]string{},
					),
				),
			},
			args:    args{u: "foo123"},
			wantErr: false,
		},
		// Testcase 2: User does not exist
		{
			name: "User does not exist",
			fields: fields{
				Clientset: fake.NewSimpleClientset(
					NewNamespace(
						"foo123",
						map[string]string{},
						map[string]string{},
					),
				),
			},
			args:    args{u: "bar113"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Clientset: tt.fields.Clientset,
			}
			if err := c.DeleteUser(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
