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
	var footer [][]string

	if err != nil {
		return err
	}

	userList, err := user.PopulateList(client, listResources)

	if err != nil {
		return err
	}

	if listResources {
		footer, err = makeFooter(userList, client)

		if err != nil {
			return err
		}
	}

	err = renderUsers(userList, footer)

	return err
}

func renderUsers(userList []user.User, footer [][]string) error {
	headerList := []string{
		"Username",
		"Full name",
		"E-mail",
		"User type",
		"Status",
	}

	if listResources {
		headerList = append(headerList, "GPU")
		headerList = append(headerList, "Mem/GPU")
		headerList = append(headerList, "Storage")
	}

	userTable, err := user.ListToTable(userList, listResources)

	if err != nil {
		return err
	}

	if listResources {
		userTable = append(userTable, footer...)
	}

	cli.RenderTable(headerList, userTable)

	return nil
}

func makeFooter(userList []user.User, client k8s.ResourceClient) ([][]string, error) {
	totalGPUs, err := client.GetTotalGPUs()
	if err != nil {
		return nil, err
	}

	resourcesUsed, err := user.TotalResourcesUsed(userList)
	if err != nil {
		return nil, err
	}

	footerList := [][]string{
		{
			"",
			"",
			"",
			"",
			"------",
			"-----",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"",
			"Total:",
			fmt.Sprint(resourcesUsed["GPU"]) + "/" + fmt.Sprint(totalGPUs),
			"",
			"",
		},
	}

	return footerList, nil
}
