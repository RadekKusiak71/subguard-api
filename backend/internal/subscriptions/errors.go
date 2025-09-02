package subscriptions

import (
	"errors"
	"net/http"

	errorx "github.com/RadekKusiak71/subguard-api/internal/errors"
)

func SubscriptionExists() errorx.APIError {
	return errorx.NewApiError(
		http.StatusConflict,
		errors.New("subscription already exist"),
	)
}

func SubscriptionNotFound() errorx.APIError {
	return errorx.NewApiError(
		http.StatusNotFound,
		errors.New("subscription not found"),
	)
}
