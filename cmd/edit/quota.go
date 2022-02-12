package edit

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uitml/quimby/internal/k8s"
	"github.com/uitml/quimby/internal/user"
	"github.com/uitml/quimby/internal/validate"
	"gopkg.in/yaml.v2"
)

func NewQuotaCmd() *cobra.Command {
	var quotaCmd = &cobra.Command{
		Use:   "res",
		Short: "Edit a users resource quota.",
		Args:  cobra.ExactArgs(1),

		RunE: RunQuota,
	}

	return quotaCmd
}

func RunQuota(cmd *cobra.Command, args []string) error {
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

	// First test: open a temporary file with VI and read the saved file.
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer tmp.Close()

	_, err = tmp.Write(y)
	if err != nil {
		return err
	}

	// Open the file in VI and read the result
	command := exec.Command("vi", tmp.Name())
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	err = command.Run()
	if err != nil {
		return err
	}

	// Process the file and apply the values
	r, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(r, md)
	if err != nil {
		return err
	}

	err = client.ApplyMetaData(username, md.Fullname, md.Email, md.Usertype)

	return err
}
