package server

import (
	"log"
	"net/http"

	"github.com/RadekKusiak71/subguard-api/internal/authentication"
	"github.com/RadekKusiak71/subguard-api/internal/errors"
	"github.com/RadekKusiak71/subguard-api/internal/users"
	"github.com/RadekKusiak71/subguard-api/internal/utils"
)

func (s *APIServer) RegisterRoutes() {
	// Stores
	userStore := users.NewStore(s.Database)

	// Handler
	authHandler := authentication.NewHandler(userStore)

	// Auth routes
	authRouter := s.Router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", MakeHandleFunc(authHandler.Register)).Methods("POST")
	authRouter.HandleFunc("/login", MakeHandleFunc(authHandler.Login)).Methods("POST")

}

type APIHandler func(w http.ResponseWriter, r *http.Request) error

func MakeHandleFunc(f APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := f(w, r)

		if err != nil {

			apiErr, ok := err.(errors.APIError)

			if ok {
				utils.WriteJSON(w, apiErr.StatusCode, apiErr)
				return
			}

			log.Printf("HandlerError: %v", err.Error())

			utils.WriteJSON(
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
