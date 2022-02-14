package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/validate"
)

// listCmd represents the list command
func newDeleteCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "rm",
		Short: "Delete user",
		Args:  cobra.ExactArgs(1),

		RunE: RunDelete,
	}

	return deleteCmd
}

func RunDelete(cmd *cobra.Command, args []string) error {
	user := args[0]

	// Validate input
	if !validate.Username(user) {
		return errors.Errorf("invalid username: %s", user)
	}

	client, err := k8s.NewClient()
	if err != nil {
		return err
	}

	v, err := client.UserExists(user)
	if err != nil {
		return err
	}
	if !v {
		return errors.Errorf("user %s does not exist", user)
	}

	// Need confirmation
	c, err := cli.Confirmation("Do you really want to delete user "+user+"? This action is irreversible.", false)
	if err != nil {
		return err
	}
	if !c {
		fmt.Printf("User %s not deleted.\n", user)
		return nil
	}

	// Do the dirty work and pray...
	err = client.DeleteUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("User %s successfully deleted. Persistent volumes must be removed manually.\n", user)

	return nil
}
