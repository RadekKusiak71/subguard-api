package authentication

import (
	"github.com/RadekKusiak71/subguard-api/internal/users"
)

func (u *RegisterUser) Validate() map[string]string {
	errors := make(map[string]string)

	if err := users.ValidateEmail(u.Email); err != nil {
		errors["email"] = err.Error()
	}

	if err := users.ValidatePassword(u.Password); err != nil {
		errors["password"] = err.Error()
	}

	if u.Password != u.PasswordConfirmation {
		errors["password_confirmation"] = "password and password confirmation does not match"
	}

	return errors
}
