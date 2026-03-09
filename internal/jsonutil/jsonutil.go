package jsonutil

import "strings"

// IsJSONLike checks if a string appears to be JSON (object or array).
func IsJSONLike(value string) bool {
	value = strings.TrimSpace(value)
	return (strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")) ||
		(strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]"))
}
