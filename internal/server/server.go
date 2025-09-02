package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type APIServer struct {
	Addr         string
	Database     *sql.DB
	Router       *mux.Router
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func NewAPIServer(addr string, database *sql.DB, writeTimeout, readTimeout time.Duration) *APIServer {
	return &APIServer{
		Addr:         addr,
		Database:     database,
		Router:       mux.NewRouter(),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
}

func (s *APIServer) Run() {
	s.RegisterRoutesAndCron()

	server := &http.Server{
		Addr:         s.Addr,
		Handler:      s.Router,
		WriteTimeout: s.WriteTimeout,
		ReadTimeout:  s.ReadTimeout,
	}

	log.Printf("Starting HTTP server on %s", s.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Couldn't start API server: %s", err.Error())
	}

}
