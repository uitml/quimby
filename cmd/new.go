/*
TODO: Implement test
*/

package cmd

import (
	"fmt"

	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/user"

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
	conf, err := cli.ParseConfig()
	if err != nil {
		return err
	}

	body, err := user.GetDefaultConfig(conf.GithubRepo, conf.GithubConfigDir+"/default-user.yaml", conf.GithubUser, conf.GithubToken)
	if err != nil {
		return err
	}

	fmt.Print(string(body))

	return nil
}
