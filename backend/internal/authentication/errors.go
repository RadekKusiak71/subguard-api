package authentication

import (
	"errors"
	"net/http"

	errorx "github.com/RadekKusiak71/subguard-api/internal/errors"
)

func InvalidCredentials() errorx.APIError {
	return errorx.NewApiError(
		http.StatusUnauthorized,
		errors.New("invalid credentials"),
	)
}
