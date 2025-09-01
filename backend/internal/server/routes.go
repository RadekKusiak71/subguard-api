package server

import (
	"github.com/RadekKusiak71/subguard-api/internal/authentication"
	"github.com/RadekKusiak71/subguard-api/internal/middlewares"
	"github.com/RadekKusiak71/subguard-api/internal/users"
	"github.com/RadekKusiak71/subguard-api/internal/utils"
)

func (s *APIServer) RegisterRoutes() {

	s.Router.Use(middlewares.LoggingMiddleware)

	// Stores
	userStore := users.NewStore(s.Database)

	// Handler
	authHandler := authentication.NewHandler(userStore)

	// Auth routes
	authRouter := s.Router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", utils.MakeHandleFunc(authHandler.Register)).Methods("POST")
	authRouter.HandleFunc("/login", utils.MakeHandleFunc(authHandler.Login)).Methods("POST")
}
