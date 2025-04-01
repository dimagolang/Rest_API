package service

import (
	"context"
	"github.com/jackc/pgx/v4"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Flight struct {
	FlightID        int    `json:"flight_id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
	DeleteAt        int64  `json:"delete_at"`
}

type FlightService struct {
	db *pgx.Conn
}

func NewFlightService(db *pgx.Conn) *FlightService {
	return &FlightService{db: db}
}

func (s *FlightService) CreateFlight(c *gin.Context) {
	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := s.db.QueryRow(context.Background(),
		"INSERT INTO flights (destination_from, destination_to) VALUES ($1, $2) RETURNING id",
		flight.DestinationFrom, flight.DestinationTo).Scan(&flight.FlightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert flight"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

func (s *FlightService) GetFlight(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight Flight
	err = s.db.QueryRow(context.Background(),
		"SELECT id, destination_from, destination_to FROM flights WHERE id=$1 AND deleted_at IS NULL", id).
		Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight details", "flight": flight})
}

func (s *FlightService) GetFlights(c *gin.Context) {
	rows, err := s.db.Query(context.Background(), "SELECT id, destination_from, destination_to FROM flights WHERE deleted_at IS NULL")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}
	defer rows.Close()

	var flights []Flight
	for rows.Next() {
		var f Flight
		if err := rows.Scan(&f.FlightID, &f.DestinationFrom, &f.DestinationTo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning flight"})
			return
		}
		flights = append(flights, f)
	}

	c.JSON(http.StatusOK, gin.H{"message": "All flights", "flights": flights})
}

func (s *FlightService) UpdateFlight(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	cmd, err := s.db.Exec(context.Background(),
		"UPDATE flights SET destination_from=$1, destination_to=$2 WHERE id=$3 AND deleted_at IS NULL",
		flight.DestinationFrom, flight.DestinationTo, id)
	if err != nil || cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found or not updated"})
		return
	}

	flight.FlightID = id
	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": flight})
}

func (s *FlightService) DeleteFlight(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	cmd, err := s.db.Exec(context.Background(),
		"UPDATE flights SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil || cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found or already deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight deleted", "id": id})
}
