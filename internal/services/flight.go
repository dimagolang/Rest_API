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
type FlightService struct {
	flights map[int]Flight
}

var (
	flightIDCounter = 1
	mutex           sync.Mutex
)

func NewFlightService() *FlightService {

	return &FlightService{
		flights: make(map[int]Flight),
	}
}

func (s *FlightService) CreateFlight(c *gin.Context) {
	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	mutex.Lock()
	flight.FlightID = flightIDCounter
	flightIDCounter++
	mutex.Unlock()

	s.flights[flight.FlightID] = flight

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

func (s *FlightService) GetFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	flight, exists := s.flights[flightID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight details", "flight": flight})
}
func (s *FlightService) GetFlights(c *gin.Context) {

	capacity := len(s.flights)

	flights := make([]Flight, 0, capacity)

	for _, flight := range s.flights {
		flights = append(flights, flight)
	}

	c.JSON(http.StatusOK, gin.H{"message": "All flights", "flights": flights})
}

func (s *FlightService) UpdateFlight(c *gin.Context) {
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

	if _, exists := s.flights[flightID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	updatedFlight.FlightID = flightID
	s.flights[flightID] = updatedFlight

	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": updatedFlight})
}

func (s *FlightService) DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	if _, exists := s.flights[flightID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	delete(s.flights, flightID)
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
