package server

import (
	"github.com/RadekKusiak71/subguard-api/internal/authentication"
	"github.com/RadekKusiak71/subguard-api/internal/middlewares"
	"github.com/RadekKusiak71/subguard-api/internal/subscriptions"
	"github.com/RadekKusiak71/subguard-api/internal/tasks"
	"github.com/RadekKusiak71/subguard-api/internal/users"
	"github.com/RadekKusiak71/subguard-api/internal/utils"
)

func (s *APIServer) RegisterRoutesAndCron() {
	// Middlewares
	s.Router.Use(middlewares.LoggingMiddleware)

	// Stores
	userStore := users.NewStore(s.Database)
	subStore := subscriptions.NewStore(s.Database)

	// Handler
	authHandler := authentication.NewHandler(userStore)
	subHandler := subscriptions.NewHandler(subStore)

	// Subrouters
	authRouter := s.Router.PathPrefix("/api/auth").Subrouter()
	subRouter := s.Router.PathPrefix("/api/subscriptions").Subrouter()

	// Authentication routes
	authRouter.HandleFunc("/register", utils.MakeHandleFunc(authHandler.Register)).Methods("POST")
	authRouter.HandleFunc("/login", utils.MakeHandleFunc(authHandler.Login)).Methods("POST")

	// Subscription routes
	subRouter.HandleFunc(
		"", utils.MakeHandleFunc(middlewares.AuthMiddleware(subHandler.ListSubscriptions, userStore)),
	).Methods("GET")

	subRouter.HandleFunc(
		"/{subscriptionID}", utils.MakeHandleFunc(middlewares.AuthMiddleware(subHandler.GetSubscription, userStore)),
	).Methods("GET")

	subRouter.HandleFunc(
		"/{subscriptionID}", utils.MakeHandleFunc(middlewares.AuthMiddleware(subHandler.UpdateSubscription, userStore)),
	).Methods("PATCH")

	subRouter.HandleFunc(
		"", utils.MakeHandleFunc(middlewares.AuthMiddleware(subHandler.CreateSubscription, userStore)),
	).Methods("POST")

	subRouter.HandleFunc(
		"/{subscriptionID}", utils.MakeHandleFunc(middlewares.AuthMiddleware(subHandler.DeleteSubscription, userStore)),
	).Methods("DELETE")

	// Cron
	subCron := tasks.NewSubscriptionCron(subStore, userStore)
	tasks.StartCron(subCron)
}
