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

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	initDB()
	defer db.Close() // Закрываем соединение при выходе

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("🚀 Server is running on port 8083...")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
