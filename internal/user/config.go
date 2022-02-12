package user

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/uitml/quimby/internal/user/reader"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Username     string `json:"username,omitempty"`
	Metadata     `json:",inline,omitempty"`
	ResourceSpec `json:",inline,omitempty"`
}

type ResourceSpec struct {
	GPU                    int `json:"gpu,inline,omitempty"`
	GPUPerJob              int `json:"gpuperjob,inline,omitempty"`
	MemoryPerJob           int `json:"memoryperjob,inline,omitempty"`
	CPUPerJob              int `json:"cpuperjob,inline,omitempty"`
	StorageProxyCPURequest int `json:"storageproxycpurequest,inline,omitempty"`
	StorageProxyCPULimit   int `json:"storageproxycpulimit,inline,omitempty"`
	StorageProxyMemory     int `json:"storageproxymemory,inline,omitempty"`
	StorageSize            int `json:"storagesize,inline,omitempty"`
}

type Metadata struct {
	Fullname string `json:"fullname,omitempty,inline"`
	Email    string `json:"email,omitempty,inline"`
	Usertype string `json:"usertype,omitempty,inline"`
}

// Populates usr given a path to a yaml file using the Reader
func (usr *Config) Populate(path string, rdr reader.Config) error {
	body, err := rdr.Read(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(body, usr)
	if err != nil {
		return err
	}
	return nil
}

// Generate config from the template in path. Populate with values from usr.
func GenerateConfig(path string, rdr reader.Config, usr Config) ([]byte, error) {
	body, err := rdr.Read(path)
	if err != nil {
		return nil, err
	}

	templ, err := template.New("default").Funcs(sprig.TxtFuncMap()).Parse(string(body))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = templ.Execute(&b, usr)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
