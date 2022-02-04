package util

func DefaultIfEmpty(field string, default_value string) string {
	if len(field) > 0 {
		return field
	}

	return default_value
}
