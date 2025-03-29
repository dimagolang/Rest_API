package services

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Flight struct {
	FlightID        int    `json:"flight_id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
}

var (
	flightIDCounter = 1
	mutex           sync.Mutex
	flights         = make(map[int]Flight)
)

func CreateFlight(c *gin.Context) {
	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	mutex.Lock()
	flight.FlightID = flightIDCounter
	flightIDCounter++
	mutex.Unlock()

	flights[flight.FlightID] = flight

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

func GetFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	flight, exists := flights[flightID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight details", "flight": flight})
}

func UpdateFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var updatedFlight Flight
	if err := c.ShouldBindJSON(&updatedFlight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if _, exists := flights[flightID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	updatedFlight.FlightID = flightID
	flights[flightID] = updatedFlight

	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": updatedFlight})
}

func DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	if _, exists := flights[flightID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	delete(flights, flightID)
	c.JSON(http.StatusOK, gin.H{"message": "Flight deleted", "id": flightID})
}

func convertID(id string) int {
	var flightID int
	_, err := fmt.Sscanf(id, "%d", &flightID)
	if err != nil {
		return 0
	}
	return flightID
}
