/*
TODO: Implement tests
*/

package usertools

import (
	"regexp"

	"github.com/uitml/quimby/internal/util"

	corev1 "k8s.io/api/core/v1"
)

type SpringfieldUser struct {
	Username string
	Fullname string
	Email    string
	Usertype string
	Status   string
}

func UserFromNamespace(namespace corev1.Namespace) SpringfieldUser {
	// I am envisioning storing resource allowance needed (e.g. memory per job) as annotations in namespace and
	// default values for resources in the cluster somehow (annotation on Springfield?).
	// Then this could be polled and populated in the list for users with default values (empty annotation)

	// TODO: Access #GPUs, storage space and memory
	usr := SpringfieldUser{
		Username: namespace.Name,
		Fullname: namespace.Annotations["springfield.uit.no/user-fullname"],
		Email:    util.DefaultIfEmpty(namespace.Annotations["springfield.uit.no/user-email"], namespace.Name+"@post.uit.no"),
		Usertype: namespace.Labels["springfield.uit.no/user-type"],
		Status:   string(namespace.Status.Phase),
	}

	return usr
}

func isValidUser(username string) bool {
	var validUser = regexp.MustCompile("^[a-z]{3}[0-9]{3}$")

	return validUser.MatchString(username)
}

func PopulateUserList(namespaceList *corev1.NamespaceList) []SpringfieldUser {
	var userList []SpringfieldUser

	for _, namespace := range namespaceList.Items {
		if isValidUser(namespace.Name) {
			userList = append(userList, UserFromNamespace(namespace))
		}
	}

	return userList
}
