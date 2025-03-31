package main

import (
	"Rest_API/internal/services"
	"Rest_API/server/http_server"
)

func main() {
	flightsService := services.NewFlightService()
	server := http_server.NewServer(flightsService)
	server.Run()
}
