package utils

import (
	"regexp"
)

// isUsernameValid returns true if the username conforms to allowed pattern:
// starts with at least one letter, and can only contain letters, digits and the following symbols ".-_".
// An arbitrary limit of 1000 characters is imposed to prevent excessively long usernames, which could cause performance issues or be used for malicious purposes.
func IsUsernameValid(username string) bool {
	re := regexp.MustCompile(`^([A-Za-z][A-Za-z0-9._-]{0,999})$`)
	return re.MatchString(username)
}

// isPasswordValid returns true if the password length is within configured bounds (12-72).
// The lower bound is to ensure a somewhat strong password.
// The upper bound is because the bcrypt algoritm can't process strings with more than 72 chars.
func IsPasswordValid(password string) bool {
	return len(password) >= 12 && len(password) <= 72
}
