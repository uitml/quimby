/*
TODO: Implement test
*/

package cmd

import (
	"fmt"

	"github.com/uitml/quimby/internal/config"
	"github.com/uitml/quimby/internal/config/reader"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/validate"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func newCreateCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new Springfield user.",

		RunE: RunGetDefault,
	}

	return createCmd
}

func RunCreate(cmd *cobra.Command, args []string) error {
	client, err := k8s.NewClient()

	if err != nil {
		return err
	}

	err = client.NewSimpleUser(args[0])

	return err
}

func RunGetDefault(cmd *cobra.Command, args []string) error {
	username := args[0]
	if !validate.Username(username) {
		return fmt.Errorf("invalid username: %s", username)
	}

	conf, err := config.Parse()
	if err != nil {
		return err
	}

	// Get default values (on github please)
	rdr := reader.Github{
		Username: conf.GithubUser,
		Token:    conf.GithubToken,
		Repo:     conf.GithubRepo,
	}
	usrConf := config.User{Username: username}
	err = usrConf.DefaultValues(conf.GithubValueDir+"/default-user.yaml", &rdr)
	if err != nil {
		return err
	}
	fmt.Println(usrConf)

	// Generate k8s user config from template
	k8sUser, err := config.GenerateConfig(conf.GithubConfigDir+"/default-user-quimby.yaml", &rdr, usrConf)
	if err != nil {
		return err
	}
	fmt.Println(string(k8sUser))

	return nil
}