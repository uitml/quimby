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

// listCmd represents the list command
func newListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List all springfield users",

		RunE: Run,
	}

	//flag = listCmd.Flags().BoolVarP("")

	return listCmd
}

func Run(cmd *cobra.Command, args []string) error {
	client, err := k8s.NewClient()

	if err != nil {
		return err
	}

	userList, err := user.PopulateList(client)

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
		"GPU",
	}

	userTable := user.ListToTable(userList)

	cli.RenderTable(headerList, userTable)
}
