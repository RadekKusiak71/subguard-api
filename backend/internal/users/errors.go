package users

import (
	"net/http"

	"github.com/RadekKusiak71/subguard-api/internal/errors"
)

func UserAlreadyExists() errors.APIError {
	return errors.NewApiError(
		http.StatusConflict,
		map[string]string{"user": "user already exists"},
	)
}
