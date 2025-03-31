package http_server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Server структура сервера с роутингом и зависимостями
type Server struct {
	flightsService flights
	cfg            config // Добавляем конфигурацию
}
type config struct {
	serverPort string
}

// NewServer создает экземпляр HTTP-сервера с настройкой роутинга
func NewServer(
	flightsService flights,
	serverPort string, // Передаем конфигурацию
) *Server {
	server := &Server{
		flightsService: flightsService,
		cfg:            config{serverPort}, // Сохраняем конфигурацию
	}

	return server
}

// Run запускает сервер
func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/flight", s.flightsService.CreateFlight)
	r.GET("/flights/:id", s.flightsService.GetFlight)
	r.GET("/flights/all", s.flightsService.GetFlights)
	r.PUT("/flights/:id", s.flightsService.UpdateFlight)
	r.DELETE("/flights/:id", s.flightsService.DeleteFlight)

	log.Printf("Server is running on port %s...", s.cfg.serverPort) // используем порт из конфигурации
	if err := r.Run(":" + s.cfg.serverPort); err != nil {           // запускаем сервер на порту из конфигурации
		log.Fatal("Failed to start server:", err)
	}
}
