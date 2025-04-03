package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log/slog"
	"os"

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
