package http_server

import "github.com/gin-gonic/gin"

type flights interface {
	CreateFlight(c *gin.Context)
	GetFlight(c *gin.Context)
	GetFlights(c *gin.Context)
	UpdateFlight(c *gin.Context)
	DeleteFlight(c *gin.Context)
	GetFlightsByCity(c *gin.Context)
}
