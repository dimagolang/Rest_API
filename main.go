package main

import (
	"Rest_API/config"
	"Rest_API/internal/services"
	"Rest_API/server/http_server"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	flightsService := services.NewFlightService()
	server := http_server.NewServer(flightsService, cfg.ServerPort)
	server.Run()
}
