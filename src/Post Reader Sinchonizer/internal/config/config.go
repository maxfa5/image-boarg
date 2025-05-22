package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	FirstConsumer Consumer
	// HTTPServer
	DataBase
}
type Consumer struct {
	Brokers string `env:"BROKERS_ADDRESSES" env-required:"true"`
	GroupId string `env:"CONSUMER_GROUP" env-required:"true"`
	Topic   string `env:"MESSAGE_TOPIC" env-required:"true"`
}
type DataBase struct {
	Username string `env:"DB_USERNAME" env-required:"true" env-default:"myuser"`
	Password string `env:"DB_PASSWORD"  env-default:"mypassword"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5430"`
	DBName   string `env:"DB_NAME" env-default:"mydatabase"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"require"`
}

func EnvLoad(logger *slog.Logger) *Config {
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("failed to load environment file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logger.Warn("CONFIG_PATH is not set")
	}

	cfg := EnvLoadInPath(configPath, logger)
	return cfg
}
func EnvLoadDb() (*DataBase, error) {

	var dbConfig DataBase

	if err := cleanenv.ReadEnv(&dbConfig); err != nil {
		return &dbConfig, fmt.Errorf("failed to load database configuration from environment: %w", err)
	}
	// fmt.Println(dbConfig)

	return &dbConfig, nil

}
func EnvLoadInPath(configPath string, logger *slog.Logger) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Warn("config file not found: ", slog.Any("error", err))
	}

	var cfg Config

	// if err := cleanenv.ReadConfig(configPath, &cfg.HTTPServer); err != nil { // Загружаем только HTTPServer
	// 	logger.Warn("config file not found", slog.Any("error", err))
	// }
	if err := cleanenv.ReadEnv(&cfg.FirstConsumer); err != nil {
		logger.Error("failed to load consumer configuration from environment:", slog.Any("error", err))
	}

	var db *DataBase
	db, err := EnvLoadDb()
	if err != nil {
		logger.Error("err db", slog.Any("error", err))
	}
	cfg.DataBase = *db
	logger.Error("DataBase configuration:", slog.Any("error", err))
	return &cfg
}
