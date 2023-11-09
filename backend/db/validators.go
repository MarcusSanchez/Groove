package db

import (
	"regexp"
)

func ValidateUser(username, email, password string) string {
	if !regexp.MustCompile(`^.{8,32}$`).MatchString(password) {
		return "invalid password: must be between 8 and 32 characters"
	}
	if !regexp.MustCompile(`^[a-zA-Z].*\d|.*[a-zA-Z].*\d|.*[a-zA-Z]\d.*$`).MatchString(password) {
		return "invalid password: must contain at least 1 uppercase letter, 1 lowercase letter, and 1 number"
	}
	if !regexp.MustCompile(`^.{4,16}$`).MatchString(username) {
		return "invalid username: must be between 4 and 16 characters"
	}
	if !regexp.MustCompile(`^.{4,320}$`).MatchString(email) {
		return "invalid email: must be between 4 and 320 characters"
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`).MatchString(email) {
		return "invalid email: must be a valid email address"
	}
	return ""
}
