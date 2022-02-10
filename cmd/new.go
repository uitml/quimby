/*
TODO: Implement test
*/

package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/uitml/quimby/internal/k8s"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func newCreateCmd() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new Springfield user.",

		RunE: RunGetConfig,
	}

	// listCmd.Flags().BoolVarP(&listResources, "show-resources", "r", false, "Show resources for all users.")

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

func RunGetConfig(cmd *cobra.Command, args []string) error {
	resp, err := http.Get("https://raw.githubusercontent.com/uitml/quimby/main/internal/validate/helpers.go")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Print(string(body))

	return nil
}
