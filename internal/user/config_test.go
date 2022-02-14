package user

import (
	"reflect"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/uitml/quimby/internal/resource"
	"github.com/uitml/quimby/internal/user/reader"
)

func TestUser_Populate(t *testing.T) {
	type fields struct {
		Username string
	}
	type args struct {
		path string
		rdr  reader.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Config
		wantErr bool
	}{
		// Read data from file
		{
			name:   "Read data from file",
			fields: fields{Username: "foo123"},
			args:   args{path: "./testdata/values.yaml", rdr: &reader.File{}},
			want: Config{
				Username: "foo123",
				Spec: &resource.Spec{
					GPU:                    pointy.Int64(2),
					GPUPerJob:              pointy.Int64(1),
					MaxMemoryPerJob:        pointy.Int64(16),
					DefaultMemoryPerJob:    pointy.Int64(12),
					CPUPerJob:              pointy.Int64(2),
					StorageProxyCPURequest: pointy.Int64(200),
					StorageProxyCPULimit:   pointy.Int64(500),
					StorageProxyMemory:     pointy.Int64(256),
					StorageSize:            pointy.Int64(500),
				},
			},
			wantErr: false,
		},
		// Read data from invalid file
		{
			name:   "Read data from invalid file",
			fields: fields{Username: "foo123"},
			args:   args{path: "./testdata/valuess.yaml", rdr: &reader.File{}},
			want: Config{
				Username: "foo123",
				Spec:     nil,
				Metadata: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usr := &Config{
				Username: tt.fields.Username,
			}
			err := usr.Populate(tt.args.path, tt.args.rdr)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Populate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*usr, tt.want) {
				t.Errorf("User.Populate() usr = %v, wantUsr %v", *usr, tt.want)
			}
		})
	}
}

func TestGenerateConfig(t *testing.T) {
	trueData, _ := (&reader.File{}).Read("./testdata/template_true.yaml")

	type args struct {
		path string
		rdr  reader.Config
		usr  Config
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// Read from file
		{
			name: "Read from file",
			args: args{
				path: "./testdata/template.yaml",
				rdr:  &reader.File{},
				usr: Config{
					Username: "foo123",
					Spec: &resource.Spec{
						GPU:                    pointy.Int64(2),
						GPUPerJob:              pointy.Int64(1),
						MaxMemoryPerJob:        pointy.Int64(16),
						DefaultMemoryPerJob:    pointy.Int64(12),
						CPUPerJob:              pointy.Int64(2),
						StorageProxyCPURequest: pointy.Int64(200),
						StorageProxyCPULimit:   pointy.Int64(500),
						StorageProxyMemory:     pointy.Int64(256),
						StorageSize:            pointy.Int64(500),
					},
				},
			},
			want:    trueData,
			wantErr: false,
		},
		// File does not exist
		{
			name: "Read from invalid file",
			args: args{
				path: "./testdata/templatefoo.yaml",
				rdr:  &reader.File{},
				usr: Config{
					Username: "foo123",
					Spec: &resource.Spec{
						GPU:                    pointy.Int64(2),
						GPUPerJob:              pointy.Int64(1),
						MaxMemoryPerJob:        pointy.Int64(16),
						DefaultMemoryPerJob:    pointy.Int64(12),
						CPUPerJob:              pointy.Int64(2),
						StorageProxyCPURequest: pointy.Int64(200),
						StorageProxyCPULimit:   pointy.Int64(500),
						StorageProxyMemory:     pointy.Int64(256),
						StorageSize:            pointy.Int64(500),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateConfig(tt.args.path, tt.args.rdr, tt.args.usr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
