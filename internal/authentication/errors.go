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

func MissingToken() errorx.APIError {
	return errorx.NewApiError(
		http.StatusUnauthorized,
		errors.New("Missing token."),
	)
}

func MissingAuthorizationHeader() errorx.APIError {
	return errorx.NewApiError(
		http.StatusUnauthorized,
		errors.New("missing authorization header"),
	)
}

func InvalidAuthorizationHeader() errorx.APIError {
	return errorx.NewApiError(
		http.StatusUnauthorized,
		errors.New("authorization header must start with 'Bearer <token>'"),
	)
}

func InvalidToken() errorx.APIError {
	return errorx.NewApiError(
		http.StatusUnauthorized,
		errors.New("invalid authorization token"),
	)
}
