/*
TODO: Implement test
*/

package cmd

import (
	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/user"

	"github.com/spf13/cobra"
)

var listResources bool

// listCmd represents the list command
func newListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List all Springfield users.",

		RunE: Run,
	}

	listCmd.Flags().BoolVarP(&listResources, "show-resources", "r", false, "Show resources for all users.")

	return listCmd
}

func Run(cmd *cobra.Command, args []string) error {
	client, err := k8s.NewClient()

	if err != nil {
		return err
	}

	userList, err := user.PopulateList(client, listResources)

	if err != nil {
		return err
	}

	renderUsers(userList)

	return nil
}

func renderUsers(userList []user.User) {
	headerList := []string{
		"Username",
		"Full name",
		"E-mail",
		"User type",
		"Status",
	}

	if listResources {
		headerList = append(headerList, "GPU")
	}

	userTable := user.ListToTable(userList, listResources)

	cli.RenderTable(headerList, userTable)
}
