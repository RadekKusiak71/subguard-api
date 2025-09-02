package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RadekKusiak71/subguard-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New() *sql.DB {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Config.DBHost,
			config.Config.DBPort,
			config.Config.DBUser,
			config.Config.DBPassword,
			config.Config.DBName,
		),
	)

	if err != nil {
		log.Fatalf("Couldn't open database: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Couldn't establish connection with database: %s", err.Error())
	}

	return db
}
