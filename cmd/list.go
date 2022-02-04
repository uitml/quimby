/*
TODO: Implement test
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/uitml/quimby/internal/usertools"

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

	userList := usertools.PopulateUserList(namespaceList)

	renderUsers(userList)

	return nil
}

func renderUsers(userList []usertools.SpringfieldUser) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Username\tFull name\tE-mail\tUser type\tStatus")
	fmt.Fprintln(w, "--------\t---------\t------\t---------\t------")

	for _, user := range userList {
		fmt.Fprintln(w, user.Username+"\t"+user.Fullname+"\t"+user.Email+"\t"+user.Usertype+"\t"+user.Status)
	}
}
