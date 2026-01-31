package utils

import (
	"regexp"
)

// isUsernameValid returns true if the username conforms to allowed pattern:
// starts with at least one letter, and can only contain letters, digits and the following symbols ".-_".
func IsUsernameValid(username string) bool {
	re := regexp.MustCompile(`^([A-Za-z][A-Za-z0-9._-]*)$`)
	return re.MatchString(username)
}

// isPasswordValid returns true if the password length is within configured bounds (12-72).
func IsPasswordValid(password string) bool {
	return len(password) >= 12 && len(password) <= 72
}
