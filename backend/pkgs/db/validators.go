package db

import (
	"errors"
	"regexp"
)

type UserField int

const (
	FieldUsername UserField = iota + 1
	FieldEmail
	FieldPassword
)

type validator struct {
	Field   UserField
	Regex   *regexp.Regexp
	Message string
}

var userValidators = []validator{
	{
		Field:   FieldPassword,
		Regex:   regexp.MustCompile(`^.{8,32}$`),
		Message: "invalid password: must be between 8 and 32 characters",
	}, {
		Field:   FieldPassword,
		Regex:   regexp.MustCompile(`^[a-zA-Z].*\d|.*[a-zA-Z].*\d|.*[a-zA-Z]\d.*$`),
		Message: "invalid password: must contain at least 1 uppercase letter, 1 lowercase letter, and 1 number",
	}, {
		Field:   FieldUsername,
		Regex:   regexp.MustCompile(`^.{4,16}$`),
		Message: "invalid username: must be between 4 and 16 characters",
	}, {
		Field:   FieldEmail,
		Regex:   regexp.MustCompile(`^.{4,320}$`),
		Message: "invalid email: must be between 4 and 320 characters",
	}, {
		Field:   FieldEmail,
		Regex:   regexp.MustCompile(`^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`),
		Message: "invalid email: must be a valid email address",
	},
}

func ValidateUser(username, email, password string) error {
	for _, v := range userValidators {
		var field string

		switch v.Field {
		case FieldUsername:
			field = username
		case FieldEmail:
			field = email
		case FieldPassword:
			field = password
		}

		if !v.Regex.MatchString(field) {
			return errors.New(v.Message)
		}
	}
	return nil
}
