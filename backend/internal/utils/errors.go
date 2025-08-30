package utils

import (
	"net/http"

	"github.com/RadekKusiak71/subguard-api/internal/errors"
)

func InvalidJSON() errors.APIError {
	return errors.NewApiError(
		http.StatusBadRequest,
		http.StatusText(http.StatusBadRequest),
	)
}

func InvalidRequest(errs map[string]string) errors.APIError {
	return errors.NewApiError(
		http.StatusUnprocessableEntity,
		errs,
	)
}
