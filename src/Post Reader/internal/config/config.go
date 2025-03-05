package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerElastic
	DataBase
}
type ServerElastic struct {
	Host        string `env:"Host_elastic" env-required:"true" env-default:"localhost"`
	Port        string `env:"Port_elastic" env-required:"true" env-default:"8085"`
	Port_server string `env:"Port_server" env-required:"true" env-default:"8085"`
	Index       string `env:"Index_elastic" env-required:"true" env-default:"masseges"`
}
type DataBase struct {
	Username string `env:"DB_USERNAME" env-required:"true" env-default:"myuser"`
	Password string `env:"DB_PASSWORD"  env-default:"mypassword"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5430"`
	DBName   string `env:"DB_NAME" env-default:"mydatabase"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"require"`
}
type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Protocol    string        `yaml:"protocol" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func EnvLoad(logger *slog.Logger) *Config {
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("failed to load environment file")
	}
	cfg := EnvLoadInPath(logger)
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

func EnvLoadInPath(logger *slog.Logger) *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg.ServerElastic); err != nil {
		logger.Error("failed to load database configuration from environment:")
		return nil
	}
	return &cfg
}
