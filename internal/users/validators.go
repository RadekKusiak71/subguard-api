package users

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) < 6 {
		return errors.New("email must have at least 6 characters")
	}

	if len(email) > 255 {
		return errors.New("email can have maximum up to 255 characters")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

func ValidatePassword(password string) error {
	password = strings.TrimSpace(password)

	if len(password) < 12 {
		return errors.New("password must have at least 12 characters")
	}

	if len(password) > 120 {
		return errors.New("password can have maximum up to 120 characters")
	}

	var hasDigit, hasUpper, hasLower bool

	for _, c := range password {
		switch {
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		}
	}

	if !hasDigit {
		return errors.New("password must have at least one digit")
	}

	if !hasUpper {
		return errors.New("password must have at least one upper case letter")
	}

	if !hasLower {
		return errors.New("password must have at least one lowercase letter")
	}

	return nil
}
