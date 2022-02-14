package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uitml/quimby/cmd/edit"
)

func newEditCmd() *cobra.Command {
	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit a user specification.",

		RunE: func(*cobra.Command, []string) error { return fmt.Errorf("missing type argument") },
	}
	editCmd.AddCommand(edit.NewMetaCmd())
	editCmd.AddCommand(edit.NewQuotaCmd())

	return editCmd
}
