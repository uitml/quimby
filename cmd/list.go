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

	return listCmd
}

func Run(cmd *cobra.Command, args []string) error {
	namespaceList := k8s.GetNamespaceList(k8s.GetClientset())
	userList := user.PopulateList(namespaceList)

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

	userTable := user.ListToTable(userList)

	cli.RenderTable(headerList, userTable)
}
