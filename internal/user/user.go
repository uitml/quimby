/*
This package implements tools and data structures for operating on users.
*/

package user

import (
	"errors"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/uitml/quimby/internal/k8s"
	internalvalidate "github.com/uitml/quimby/internal/validate"

	corev1 "k8s.io/api/core/v1"
)

type User struct {
	Username      string
	fullname      string
	email         string
	usertype      string
	ResourceQuota k8s.ResourceQuota
}

func FromNamespace(namespace corev1.Namespace) User {
	// I am envisioning storing resource allowance needed (e.g. memory per job) as annotations in namespace and
	// default values for resources in the cluster somehow (annotation on Springfield?).
	// Then this could be polled and populated in the list for users with default values (empty annotation)
	usr := User{
		Username: namespace.Name,
		fullname: namespace.Annotations[k8s.AnnotationUserFullname],
		email:    internalvalidate.DefaultIfEmpty(namespace.Annotations[k8s.AnnotationUserEmail], namespace.Name+"@post.uit.no"),
		usertype: namespace.Labels[k8s.LabelUserType],
	}

	return usr
}

func PopulateList(c k8s.ResourceClient, listResources bool) ([]User, error) {
	var userList []User

	namespaceList, err := c.GetNamespaceList()

	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		if internalvalidate.Username(namespace.Name) {
			newUser := FromNamespace(namespace)

			// Will only poll for resources if flag is true (for efficiency)
			if listResources {
				newUser.ResourceQuota, err = c.GetResourceQuota(namespace.Name)
				if err != nil {
					return nil, err
				}
			}

			userList = append(userList, newUser)
		}
	}

	if len(userList) == 0 {
		return userList, errors.New("no users found on the cluster")
	}

	return userList, nil
}

func ListToTable(userList []User, listResources bool) ([][]string, error) {
	var table [][]string

	for i, usr := range userList {
		table = append(table, []string{
			usr.Username,
			usr.fullname,
			usr.email,
			usr.usertype,
		})

		// Only show resources if the user has asked for it
		if listResources {
			m := memoryPerGPU(usr)

			table[i] = append(table[i], fmt.Sprint(usr.ResourceQuota.GPU.Used)+"/"+fmt.Sprint(usr.ResourceQuota.GPU.Max))
			table[i] = append(table[i], humanize.IBytes(uint64(m)))
			table[i] = append(table[i], humanize.IBytes(uint64(usr.ResourceQuota.Storage)))
		}
	}

	return table, nil
}

func (usr *User) Metadata() *Metadata {
	return &Metadata{
		Fullname: usr.fullname,
		Email:    usr.email,
		Usertype: usr.usertype,
	}
}
