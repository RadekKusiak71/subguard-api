package utils

import (
	"errors"
	"net/http"

	errorx "github.com/RadekKusiak71/subguard-api/internal/errors"
)

func InvalidJSON() errorx.APIError {
	return errorx.NewApiError(
		http.StatusBadRequest,
		errors.New("invalid JSON reqeust data"),
	)
}

func InvalidRequest(errors map[string]string) errorx.APIError {
	return errorx.APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    errors,
	}
}
