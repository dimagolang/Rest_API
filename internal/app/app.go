package app

import (
	"Rest_API/internal/config"
	"Rest_API/internal/http_server"
	"Rest_API/internal/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var viperInstance *viper.Viper
var dbConn *pgx.Conn

type App struct {
	cfg    *config.Config
	server *http_server.Server
	db     *pgx.Conn // Добавлено: подключение к базе
}

func (a *App) Init() {
	// Инициализация логгера
	if err := InitLogger("debug", "", true); err != nil {
		log.Fatal().Err(err).Msg("Ошибка инициализации логгера")
		return
	}

	const op = "App.Init"

	// Загрузка конфигурации
	conf, err := config.GetConfigReader("")
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
		return
	}

	viperInstance = conf

	// Подключение к базе данных
	a.cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	dbConn, err = config.GetDBConnect(a.cfg)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	// Создание сервиса
	flightsService := service.NewFlightService(a.db)

	// Создание HTTP-сервера
	a.server = http_server.NewServer(flightsService, *a.cfg)

	log.Info().Msg(fmt.Sprintf("Инициализация завершена. Сервер будет запущен на порту %s", a.cfg.ServerPort))
}

func (a *App) Run() {
	a.Init()
	defer a.db.Close(context.Background())
	a.server.Run()
}
