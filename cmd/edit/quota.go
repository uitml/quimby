package edit

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/resource"
	"github.com/uitml/quimby/internal/user"
	"github.com/uitml/quimby/internal/user/reader"
	"github.com/uitml/quimby/internal/validate"
	"gopkg.in/yaml.v2"
)

func NewQuotaCmd() *cobra.Command {
	var quotaCmd = &cobra.Command{
		Use:   "quota",
		Short: "Edit a users resource quota.",
		Args:  cobra.ExactArgs(1),

		RunE: RunQuota,
	}

	return quotaCmd
}

func RunQuota(cmd *cobra.Command, args []string) error {
	username := args[0]

	// Validate input
	if !validate.Username(username) {
		return errors.Errorf("invalid username: %s", username)
	}

	client, err := k8s.NewClient()
	if err != nil {
		return err
	}

	u, err := client.UserExists(username)
	if err != nil {
		return err
	}
	if !u {
		return fmt.Errorf("user %s does not exist", username)
	}

	// Get current values
	spec, err := client.Spec(username)
	if err != nil {
		return err
	}

	// Selecting which values are allowed to be edited...
	tmpSpec := resource.Spec{
		GPU:                 spec.GPU,
		MaxMemoryPerJob:     spec.MaxMemoryPerJob,
		DefaultMemoryPerJob: spec.DefaultMemoryPerJob,
	}

	// Edit values
	s, err := yaml.Marshal(tmpSpec)
	if err != nil {
		return err
	}
	es, err := cli.Editor(s)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(es, spec)
	if err != nil {
		return err
	}

	// Populate manifests
	conf, err := cli.ParseConfig()
	if err != nil {
		return err
	}

	rdr := reader.Github{
		Username: conf.GithubUser,
		Token:    conf.GithubToken,
		Repo:     conf.GithubRepo,
	}
	usrConf := user.Config{Username: username, Spec: spec}
	k8sUser, err := user.GenerateConfig(conf.GithubConfigDir+"/default-user-quimby.yaml", &rdr, usrConf)
	if err != nil {
		return err
	}

	// Apply updates
	err = client.Apply(username, k8sUser)

	return err
}
