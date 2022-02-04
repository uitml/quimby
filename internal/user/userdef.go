/*
This package implements tools and data structures for operating on users.
*/

package user

import (
	internalvalidate "github.com/uitml/quimby/internal/validate"

	corev1 "k8s.io/api/core/v1"
)

type User struct {
	Username string
	Fullname string
	Email    string
	Usertype string
	Status   string
}

func FromNamespace(namespace corev1.Namespace) User {
	// I am envisioning storing resource allowance needed (e.g. memory per job) as annotations in namespace and
	// default values for resources in the cluster somehow (annotation on Springfield?).
	// Then this could be polled and populated in the list for users with default values (empty annotation)

	// TODO: Access #GPUs, storage space and memory
	usr := User{
		Username: namespace.Name,
		Fullname: namespace.Annotations["springfield.uit.no/user-fullname"],
		Email:    internalvalidate.DefaultIfEmpty(namespace.Annotations["springfield.uit.no/user-email"], namespace.Name+"@post.uit.no"),
		Usertype: namespace.Labels["springfield.uit.no/user-type"],
		Status:   string(namespace.Status.Phase),
	}

	return usr
}

func PopulateList(namespaceList *corev1.NamespaceList) []User {
	var userList []User

	for _, namespace := range namespaceList.Items {
		if internalvalidate.Username(namespace.Name) {
			userList = append(userList, FromNamespace(namespace))
		}
	}

	return userList
}
