package utils

import "strings"

// If the value is "no", "none", "disabled", "disable" or "false" (case insensitive), it returns false. Otherwise, it returns true.
func IsActive(value string) bool {
	switch strings.ToLower(value) {
	case "no", "none", "disabled", "disable", "false":
		return false
	default:
		return true
	}
}
