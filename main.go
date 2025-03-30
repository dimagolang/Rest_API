package main

import (
	"Rest_API/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	flightService := services.NewFlightService()
	r.POST("/flights", flightService.CreateFlight)
	r.GET("/flights/:id", flightService.GetFlight)
	r.PUT("/flights/:id", flightService.UpdateFlight)
	r.DELETE("/flights/:id", flightService.DeleteFlight)

	log.Println("Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
