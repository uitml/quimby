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

		RunE: RunList,
	}

	listCmd.Flags().BoolVarP(&listResources, "show-resources", "r", false, "Show resources for all users.")

	return listCmd
}

func RunList(cmd *cobra.Command, args []string) error {
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
	headers := [][]string{
		{
			"Username",
			"Full name",
			"E-mail",
			"User type",
		},
	}

	if listResources {
		headers[0] = append(headers[0], "GPU", "Mem/GPU", "Storage")
	}

	userTable, err := user.ListToTable(userList, listResources)
	if err != nil {
		return err
	}

	if listResources {
		cli.RenderTable(headers, userTable, footer)
	} else {
		cli.RenderTable(headers, userTable)
	}

	return nil
}

func makeFooter(userList []user.User, client k8s.ResourceClient) ([][]string, error) {
	GPUSummary, err := client.GetTotalGPUs()
	if err != nil {
		return nil, err
	}

	resourceUsage := user.TotalResourcesUsed(userList)

	footer := [][]string{
		{
			"",
			"",
			"",
			"Total:",
			fmt.Sprint(resourceUsage[k8s.ResourceGPU]) + "/" + fmt.Sprint(GPUSummary.Max),
			"",
			"",
		},
	}

	return footer, nil
}
