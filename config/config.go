package config

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config struct holds the database configuration
type Config struct {
	ServerPort string
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

		ServerPort: viper.GetString("SERVER_PORT"),
	}

	slog.Info("Configuration loaded successfully")
	return config, nil
}
