package app

import (
	"Rest_API/internal/service"
	"Rest_API/server/http_server"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var viperInstance *viper.Viper

type App struct {
	cfg    *Config
	server *http_server.Server // Добавляем поле для сервера
}

func (a *App) Init() {
	// Инициализация логгера с использованием zerolog
	if err := InitLogger("debug", "", true); err != nil {
		log.Fatal().Err(err).Msg("Ошибка инициализации логгера")
		return
	}

	const op string = "App.Init"

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
	a.server = http_server.NewServer(flightsService, a.cfg.ServerPort) // Сохраняем сервер в поле структуры
}

func (a *App) setConfig() {
	a.cfg = &Config{}
	a.cfg.ServerPort = viperInstance.GetString("SERVER_PORT") // Получаем порт из viper
}

func (a *App) Run() {
	// Инициализация приложения
	a.Init()

	// Запуск сервера
	a.server.Run() // Используем поле структуры для запуска сервера

	// Логируем успешный запуск сервера
	log.Info().Msg("Сервер успешно запущен на порту " + a.cfg.ServerPort)
}
