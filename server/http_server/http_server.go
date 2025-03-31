package http_server

import (
	"Rest_API/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Server структура сервера с роутингом и зависимостями
type Server struct {
	flightsService *services.FlightService
}

// NewServer создает экземпляр HTTP-сервера с настройкой роутинга
func NewServer(

	flightsService *services.FlightService,

) *Server {
	server := &Server{

		flightsService: flightsService,
	}

	return server
}

// запускает сервер
func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/flights", s.flightsService.CreateFlight)
	r.GET("/flights/:id", s.flightsService.GetFlight)
	r.PUT("/flights/:id", s.flightsService.UpdateFlight)
	r.DELETE("/flights/:id", s.flightsService.DeleteFlight)

	log.Println("Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
