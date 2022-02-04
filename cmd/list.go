/*
TODO: Implement test
*/

package cmd

import (
	"context"

	"github.com/uitml/quimby/internal/cli"
	"github.com/uitml/quimby/internal/user"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	// TODO: Client abstraction (this needs to be done for every command).
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)

	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	userList := user.PopulateList(namespaceList)

	renderUsers(userList)

	return nil
}

func renderUsers(userList []user.User) {
	var rowList [][]string

	headerList := []string{
		"Username",
		"Full name",
		"E-mail",
		"User type",
		"Status",
	}

	for _, user := range userList {
		rowList = append(rowList, []string{
			user.Username,
			user.Fullname,
			user.Email,
			user.Usertype,
			user.Status,
		})
	}

	cli.RenderTable(headerList, rowList)
}
