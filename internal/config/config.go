package config

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"path/filepath"

	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/spf13/viper"
)

const defaultConfigFilePath = "./.env"

// Config struct holds the database configuration
type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	DatabaseURL string
}

// LoadConfig initializes Viper and loads configuration
func LoadConfig() (*Config, error) {
	// Load .env file if present
	_ = godotenv.Load()

	viper.SetConfigName("config")   // config file name (without extension)
	viper.SetConfigType("env")      // config file type
	viper.AddConfigPath(".")        // path to look for the config file
	viper.AddConfigPath("./config") // path inside a config folder

	// Enable environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("No config file found, using system environment variables")
	}

	config := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		ServerPort: viper.GetString("SERVER_PORT"),
	}

	slog.Info("Configuration loaded successfully")
	return config, nil
}

func GetDBConnect(cfg *Config) (*pgx.Conn, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	slog.Info("Connecting to Postgres", "url", dbURL)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return nil, errors.Wrap(err, "unable to connect to database")
	}
	log.Println("Running PostgreSQL migrations")
	if err := runPgMigrations(dbURL, "migrations"); err != nil {
		return nil, errors.Wrap(err, "runPgMigrations failed")
	}
	slog.Info("Connected to PostgreSQL successfully")
	return conn, nil
}

func GetConfigReader(path string) (*viper.Viper, error) {
	const op = "app.GetConfigReader"

	if path == "" {
		path = defaultConfigFilePath
	}

	conf := viper.New()
	conf.SetConfigFile(path)
	conf.SetConfigType("env")
	conf.AutomaticEnv()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = conf.WriteConfigAs(path)
		if err != nil {
			return nil, fmt.Errorf("%s: creating config file error: %w", op, err)
		}
	}

	err := conf.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("%s: config reading error: %w", op, err)
	}

	return conf, nil
}

// runPgMigrations runs Postgres migrations
func runPgMigrations(dsn, path string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	slog.Info("Running migrations...", "dsn", dsn, "path", path)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return errors.Wrap(err, "failed to open DB connection for migrations")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create migration driver")
	}

	// Convert relative to absolute path and prepend file://
	// ✅ Получаем абсолютный путь с прямыми слэшами

	// 🔥 Преобразуем Windows-путь в корректный file:// URL
	sourceURL, err := getFileMigrationURL("./migrations")
	if err != nil {
		slog.Error("invalid migration path", "error", err)
		return err
	}

	slog.Info("Resolved migration path", "sourceURL", sourceURL)

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "migration failed")
	}

	slog.Info("Migrations applied successfully")
	return nil
}
func getFileMigrationURL(relPath string) (string, error) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", err
	}

	// Используем url.PathEscape чтобы защититься от пробелов, кириллицы и т.п.
	u := url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absPath),
	}
	return u.String(), nil
}
