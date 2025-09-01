package utils

import (
	"log"
	"net/http"

	"github.com/RadekKusiak71/subguard-api/internal/errors"
)

type APIHandler func(w http.ResponseWriter, r *http.Request) error

func MakeHandleFunc(f APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := f(w, r)

		if err != nil {

			apiErr, ok := err.(errors.APIError)

			if ok {
				WriteJSON(w, apiErr.StatusCode, apiErr)
				return
			}

			log.Printf("HandlerError: %v", err.Error())

			WriteJSON(
				w,
				http.StatusInternalServerError,
				map[string]any{
					"status_code": http.StatusInternalServerError,
					"message":     http.StatusText(http.StatusInternalServerError),
				},
			)

		}

	}
}
