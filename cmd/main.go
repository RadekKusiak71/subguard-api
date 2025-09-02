package main

import (
	"fmt"
	"time"

	"github.com/RadekKusiak71/subguard-api/internal/config"
	"github.com/RadekKusiak71/subguard-api/internal/database"
	"github.com/RadekKusiak71/subguard-api/internal/server"
)

func main() {

	server := server.NewAPIServer(
		fmt.Sprintf(":%s", config.Config.Port),
		database.New(),
		15*time.Second,
		15*time.Second,
	)

	server.Run()

}
