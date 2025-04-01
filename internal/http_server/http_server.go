package http_server

import (
	"Rest_API/internal/config"
	"Rest_API/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Server структура сервера с роутингом и зависимостями
type Server struct {
	flightsService *service.FlightService
	cfg            config.Config
}

// NewServer создает экземпляр HTTP-сервера с настройкой роутинга
func NewServer(
	flightsService *service.FlightService,
	cfg config.Config,
) *Server {
	return &Server{
		flightsService: flightsService,
		cfg:            cfg,
	}
}

// Run запускает сервер
func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Роуты для рейсов
	r.POST("/flight", s.flightsService.CreateFlight)
	r.GET("/flights/:id", s.flightsService.GetFlight)
	r.GET("/flights/all", s.flightsService.GetFlights)
	r.PUT("/flights/:id", s.flightsService.UpdateFlight)
	r.DELETE("/flights/:id", s.flightsService.DeleteFlight)

	log.Printf("Server is running on port %s...", s.cfg.ServerPort)
	if err := r.Run(":" + s.cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
