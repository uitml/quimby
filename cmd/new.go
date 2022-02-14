/*
TODO: Implement test
*/

package cmd

import (
	"fmt"

	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/user"
	"github.com/uitml/quimby/internal/user/reader"
	"github.com/uitml/quimby/internal/validate"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func newCreateCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new Springfield user.",
		Args:  cobra.ExactArgs(1),

		RunE: RunNew,
	}

	return createCmd
}

func RunNew(cmd *cobra.Command, args []string) error {
	username := args[0]
	if !validate.Username(username) {
		return fmt.Errorf("invalid username: %s", username)
	}

	conf, err := cli.ParseConfig()
	if err != nil {
		return err
	}

	// Get default values (on github please)
	rdr := reader.Github{
		Username: conf.GithubUser,
		Token:    conf.GithubToken,
		Repo:     conf.GithubRepo,
	}
	usrConf := user.Config{Username: username}
	err = usrConf.Populate(conf.GithubValueDir+"/default-user.yaml", &rdr)
	if err != nil {
		return err
	}

	// Generate k8s user config from template
	k8sUser, err := user.GenerateConfig(conf.GithubConfigDir+"/default-user-quimby.yaml", &rdr, usrConf)
	if err != nil {
		return err
	}

	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	client.Apply(username, k8sUser)

	return nil
}
