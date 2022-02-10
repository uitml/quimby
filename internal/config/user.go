package config

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/uitml/quimby/internal/config/reader"
	"gopkg.in/yaml.v2"
)

type User struct {
	Username               string `yaml:"Username,omitempty"`
	GPU                    int    `yaml:"GPU"`
	GPUPerJob              int    `yaml:"GPUPerJob"`
	MemoryPerJob           int    `yaml:"MemoryPerJob"`
	CPUPerJob              int    `yaml:"CPUPerJob"`
	StorageProxyCPURequest int    `yaml:"StorageProxyCPURequest"`
	StorageProxyCPULimit   int    `yaml:"StorageProxyCPULimit"`
	StorageProxyMemory     int    `yaml:"StorageProxyMemory"`
	StorageSize            int    `yaml:"StorageSize"`
}

// Populates usr given a path to a yaml file using the Reader
func (usr *User) DefaultValues(path string, rdr reader.Reader) error {
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
func GenerateConfig(path string, rdr reader.Reader, usr User) ([]byte, error) {
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
