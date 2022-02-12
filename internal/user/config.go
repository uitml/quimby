package user

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/uitml/quimby/internal/resource"
	"github.com/uitml/quimby/internal/user/reader"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Username       string `yaml:"username,omitempty"`
	*Metadata      `yaml:"metadata,omitempty"`
	*resource.Spec `yaml:"resourcespec,omitempty"`
}

type Metadata struct {
	Fullname string `yaml:"fullname,omitempty"`
	Email    string `yaml:"email,omitempty"`
	Usertype string `yaml:"usertype,omitempty"`
}

// Populates usr given a path to a yaml file using the Reader
func (usr *Config) Populate(path string, rdr reader.Config) error {
	body, err := rdr.Read(path)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
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
