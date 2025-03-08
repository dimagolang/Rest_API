package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error

	dsn := "postgres://postgres:postgres@localhost:5432/flights_db?sslmode=disable"

	// Подключение с 10 попытками
	for i := 1; i <= 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Attempt %d: Error opening database connection: %v\n", i, err)
		} else if err = db.Ping(); err == nil {
			log.Println("✅ Connected to PostgreSQL!")
			return
		} else {
			log.Printf("Attempt %d: Error pinging database: %v\n", i, err)
		}

		log.Println("⏳ Waiting for database... Attempt", i)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("❌ Failed to connect to database after multiple attempts:", err)
}

// ... существующий код ...

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	initDB()
	defer db.Close() // Закрываем соединение при выходе

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Добавление новых маршрутов для работы с рейсами
	r.POST("/flights", createFlight)       // Создание рейса
	r.GET("/flights/:id", getFlight)       // Получение рейса по ID
	r.PUT("/flights/:id", updateFlight)    // Обновление рейса по ID
	r.DELETE("/flights/:id", deleteFlight) // Удаление рейса по ID

	log.Println("🚀 Server is running on port 8083...")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

// Обработчик для создания рейса
func createFlight(c *gin.Context) {
	// Логика создания рейса
	c.JSON(http.StatusCreated, gin.H{"message": "Flight created"})
}

// Обработчик для получения рейса
func getFlight(c *gin.Context) {
	id := c.Param("id")
	// Логика получения рейса по ID
	c.JSON(http.StatusOK, gin.H{"message": "Flight details", "id": id})
}

// Обработчик для обновления рейса
func updateFlight(c *gin.Context) {
	id := c.Param("id")
	// Логика обновления рейса по ID
	c.JSON(http.StatusOK, gin.H{"message": "Flight updated", "id": id})
}

// Обработчик для удаления рейса
func deleteFlight(c *gin.Context) {
	id := c.Param("id")
	// Логика удаления рейса по ID
	c.JSON(http.StatusOK, gin.H{"message": "Flight deleted", "id": id})
}
