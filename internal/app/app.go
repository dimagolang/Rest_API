package app

import (
	"Rest_API/internal/service"
	"Rest_API/server/http_server"
	"github.com/rs/zerolog"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var viperInstance *viper.Viper

type App struct {
	cfg *Config
}

func (a *App) Init() {
	// Инициализация логгера с использованием zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	const logTrace string = "App.Init"

	// Загрузка конфигурации
	conf, err := GetConfigReader("")
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
		return
	}

	viperInstance = conf
	a.setConfig()

	// Создание сервиса для работы с рейсами
	flightsService := service.NewFlightService()

	// Создание HTTP-сервера
	server := http_server.NewServer(flightsService, a.cfg.ServerPort) // Используем a.cfg.ServerPort

	// Запуск сервера
	server.Run()

	log.Info().Msg("Сервер успешно запущен на порту " + a.cfg.ServerPort) // Используем a.cfg.ServerPort
}

func (a *App) setConfig() {
	a.cfg = &Config{}
	a.cfg.ServerPort = viperInstance.GetString("SERVER_PORT") // Получаем порт из viper
}

func (a *App) Run() {
	a.Init()
}
