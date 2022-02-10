package config

import (
	"reflect"
	"testing"

	"github.com/uitml/quimby/internal/config/reader"
)

func TestUser_DefaultValues(t *testing.T) {
	type fields struct {
		Username               string
		GPU                    int
		GPUPerJob              int
		MemoryPerJob           int
		CPUPerJob              int
		StorageProxyCPURequest int
		StorageProxyCPULimit   int
		StorageProxyMemory     int
		StorageSize            int
	}
	type args struct {
		path string
		rdr  reader.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    User
		wantErr bool
	}{
		// Read data from file
		{
			name:   "Read data from file",
			fields: fields{Username: "foo123"},
			args:   args{path: "./testdata/values.yaml", rdr: &reader.File{}},
			want: User{
				Username:               "foo123",
				GPU:                    2,
				GPUPerJob:              1,
				MemoryPerJob:           16,
				CPUPerJob:              2,
				StorageProxyCPURequest: 200,
				StorageProxyCPULimit:   500,
				StorageProxyMemory:     256,
				StorageSize:            500,
			},
			wantErr: false,
		},
		// Read data from invalid file
		{
			name:   "Read data from invalid file",
			fields: fields{Username: "foo123"},
			args:   args{path: "./testdata/valuess.yaml", rdr: &reader.File{}},
			want: User{
				Username:               "foo123",
				GPU:                    0,
				GPUPerJob:              0,
				MemoryPerJob:           0,
				CPUPerJob:              0,
				StorageProxyCPURequest: 0,
				StorageProxyCPULimit:   0,
				StorageProxyMemory:     0,
				StorageSize:            0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usr := &User{
				Username:               tt.fields.Username,
				GPU:                    tt.fields.GPU,
				GPUPerJob:              tt.fields.GPUPerJob,
				MemoryPerJob:           tt.fields.MemoryPerJob,
				CPUPerJob:              tt.fields.CPUPerJob,
				StorageProxyCPURequest: tt.fields.StorageProxyCPURequest,
				StorageProxyCPULimit:   tt.fields.StorageProxyCPULimit,
				StorageProxyMemory:     tt.fields.StorageProxyMemory,
				StorageSize:            tt.fields.StorageSize,
			}
			err := usr.DefaultValues(tt.args.path, tt.args.rdr)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.DefaultValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *usr != tt.want {
				t.Errorf("User.DefaultValues() usr = %v, wantUsr %v", *usr, tt.want)
			}
		})
	}
}

func TestGenerateConfig(t *testing.T) {
	trueData, _ := (&reader.File{}).Read("./testdata/template_true.yaml")

	type args struct {
		path string
		rdr  reader.Reader
		usr  User
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
				usr: User{
					Username:               "foo123",
					GPU:                    2,
					GPUPerJob:              1,
					MemoryPerJob:           16,
					CPUPerJob:              2,
					StorageProxyCPURequest: 200,
					StorageProxyCPULimit:   500,
					StorageProxyMemory:     256,
					StorageSize:            500,
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
				usr: User{
					Username:               "foo123",
					GPU:                    2,
					GPUPerJob:              1,
					MemoryPerJob:           16,
					CPUPerJob:              2,
					StorageProxyCPURequest: 200,
					StorageProxyCPULimit:   500,
					StorageProxyMemory:     256,
					StorageSize:            500,
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
