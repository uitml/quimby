package edit

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/user"
	"github.com/uitml/quimby/internal/validate"
	"gopkg.in/yaml.v2"
)

func NewMetaCmd() *cobra.Command {
	var metaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Edit a users metadata.",
		Args:  cobra.ExactArgs(1),

		RunE: RunMeta,
	}

	return metaCmd
}

func RunMeta(cmd *cobra.Command, args []string) error {
	username := args[0]

	// Validate input
	if !validate.Username(username) {
		return errors.Errorf("invalid username: %s", username)
	}

	client, err := k8s.NewClient()
	if err != nil {
		return err
	}

	ns, err := client.Namespace(username)
	if err != nil {
		return err
	}
	u := user.FromNamespace(*ns)
	md := u.Metadata()
	y, err := yaml.Marshal(md)
	if err != nil {
		return err
	}

	r, err := cli.Editor(y)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(r, md)
	if err != nil {
		return err
	}

	err = client.ApplyMetadata(username, md.Fullname, md.Email, md.Usertype)

	return err
}
