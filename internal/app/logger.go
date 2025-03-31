package app

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger инициализирует глобальный логгер с заданным уровнем и опционально лог-файлом
func InitLogger(level, logFile string, humanReadable bool) error {
	// Разбираем уровень логирования
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("ошибка инициализации логгера: %w", err)
	}
	zerolog.SetGlobalLevel(lvl)

	// Определяем writer (Stdout или файл)
	var writers []io.Writer

	if humanReadable {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		writers = append(writers, consoleWriter)
	} else {
		writers = append(writers, os.Stdout)
	}

	// Логирование в файл, если задан путь
	if strings.TrimSpace(logFile) != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("ошибка открытия лог-файла: %w", err)
		}
		writers = append(writers, file)
	}

	// Объединяем несколько выходов (Stdout + файл)
	multiWriter := zerolog.MultiLevelWriter(writers...)

	// Создаем глобальный логгер

	log.Logger = zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller(). // Добавляет caller (файл и строка вызова)
		Logger()

	return nil
}
