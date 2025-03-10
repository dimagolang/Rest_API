package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

// Структура рейса
type Flight struct {
	FlightID        int    `json:"flight_id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
}

// Глобальный счетчик ID и мьютекс для безопасного увеличения
var (
	flightIDCounter = 1
	mutex           sync.Mutex
)

// Хранилище рейсов (в памяти)
var flights = make(map[int]Flight)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/flights", createFlight)
	r.GET("/flights/:id", getFlight)
	r.PUT("/flights/:id", updateFlight)
	r.DELETE("/flights/:id", deleteFlight)

	log.Println("Server is running on port 8083...")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// ✅ Создание рейса с автоматическим увеличением flight_id
func createFlight(c *gin.Context) {
	var flight Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Блокируем мьютекс для безопасного увеличения счетчика
	mutex.Lock()
	flight.FlightID = flightIDCounter
	flightIDCounter++
	mutex.Unlock()

	// Сохраняем рейс в памяти
	flights[flight.FlightID] = flight

	c.JSON(http.StatusCreated, gin.H{"message": "Flight created", "flight": flight})
}

// ✅ Получение информации о рейсе
func getFlight(c *gin.Context) {
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

// ✅ Обновление рейса
func updateFlight(c *gin.Context) {
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

	// Обновляем данные рейса
	updatedFlight.FlightID = flightID
	flights[flightID] = updatedFlight

	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "flight": updatedFlight})
}

// ✅ Удаление рейса
func deleteFlight(c *gin.Context) {
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

// Функция для конвертации ID в int
func convertID(id string) int {
	var flightID int
	_, err := fmt.Sscanf(id, "%d", &flightID)
	if err != nil {
		return 0
	}
	return flightID
}
