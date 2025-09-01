package users

import (
	"errors"
	"net/http"

	errorx "github.com/RadekKusiak71/subguard-api/internal/errors"
)

func UserAlreadyExist() errorx.APIError {
	return errorx.NewApiError(
		http.StatusConflict,
		errors.New("user already exist"),
	)
}

func UserDoesNotExist() errorx.APIError {
	return errorx.NewApiError(
		http.StatusNotFound,
		errors.New("user does not exist"),
	)
}

func AccountNotVerified() errorx.APIError {
	return errorx.NewApiError(
		http.StatusForbidden,
		errors.New("user account is not verified"),
	)
}
