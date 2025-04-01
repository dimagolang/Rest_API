package service

import (
	"Rest_API/internal/models"
	"Rest_API/internal/repository"
	"context"
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
	flightsRepo *repository.FlightsRepo
	flights     map[int]Flight
}

var (
	flightIDCounter = 1
	mutex           sync.Mutex
)

func NewFlightService(flightsRepo *repository.FlightsRepo) *FlightService {

	return &FlightService{
		flightsRepo: flightsRepo,
		flights:     make(map[int]Flight),
	}
}

func (s *FlightService) CreateFlight(c *gin.Context) {
	var flight models.Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	flight.DeleteAt = 0

	err := s.flightsRepo.InsertFlightToDB(context.Background(), &flight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

func (s *FlightService) GetFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	flight, err := s.flightsRepo.GetFlightByIDFromDB(context.Background(), flightID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight details", "flight": flight})
}

func (s *FlightService) GetFlights(c *gin.Context) {
	flights, err := s.flightsRepo.GetAllFlightsFromDB(context.Background())
	if err != nil {
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

	var updatedFlight models.Flight
	if err := c.ShouldBindJSON(&updatedFlight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Устанавливаем ID из URL
	updatedFlight.FlightID = flightID

	err := s.flightsRepo.UpdateFlightInDB(context.Background(), &updatedFlight)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found or not updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": updatedFlight})
}

func (s *FlightService) DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	flightID := convertID(id)
	if flightID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	err := s.flightsRepo.DeleteFlightFromDB(context.Background(), flightID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found or already deleted"})
		return
	}

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
