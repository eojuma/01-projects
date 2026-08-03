package validation

import "regexp"

var (
	emailRegex    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
)

// ValidateEmail checks whether the provided value looks like a valid email.
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword enforces a minimum password length.
func ValidatePassword(password string) bool {
	return len(password) >= 8
}

// ValidateUsername enforces allowed characters and length.
func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
