package db

import "regexp"

func ValidateUser(username, email string) string {
	if regexp.MustCompile(`^.{4,16}$`).MatchString(username) {
		return "invalid username: must be between 4 and 16 characters"
	}
	if regexp.MustCompile(`^.{4,320}$`).MatchString(email) {
		return "invalid email: must be between 4 and 320 characters"
	}
	if regexp.MustCompile(`^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`).MatchString(email) {
		return "invalid email: must be a valid email address"
	}
	return "invalid username or email"
}
