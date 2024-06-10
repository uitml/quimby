package validate

import "regexp"

func DefaultIfEmpty(field string, default_value string) string {
	if len(field) > 0 {
		return field
	}

	return default_value
}

func Username(username string) bool {
	var validUsername = regexp.MustCompile("^[a-z]{3,5}[0-9]{3,4}$")

	return validUsername.MatchString(username)
}
