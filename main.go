package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

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

	log.Println("🚀 Server is running on port 8083...")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

func createFlight(c *gin.Context) {
	// Логика создания рейса
	c.JSON(http.StatusCreated, gin.H{"message": "Flight created"})
}

func getFlight(c *gin.Context) {
	id := c.Param("id")
	// Логика получения рейса по ID
	c.JSON(http.StatusOK, gin.H{"message": "Flight details", "id": id})
}

func updateFlight(c *gin.Context) {
	id := c.Param("id")
	// Логика обновления рейса по ID
	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "id": id})
}

func deleteFlight(c *gin.Context) {
	id := c.Param("id")
	// Логика удаления рейса по ID
	c.JSON(http.StatusOK, gin.H{"message": "Flight deleted", "id": id})
}
