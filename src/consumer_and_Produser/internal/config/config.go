package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	FirstConsumer Consumer `yaml:"consumer_1" env-required:"true"`
	HTTPServer    `yaml:"http_server" env-required:"true"`
	DataBase      `yaml:`
}
type Consumer struct {
	Brokers string `yaml:"env" env-required:"true"`
	GroupId string `yaml:"consumer_group" env-required:"true"`
	Topic   string `yaml:"message_topic" env-required:"true"`
}
type DataBase struct {
	Username string `envconfig:"DB_USERNAME" env-default:"myuser"`
	Password string `envconfig:"DB_PASSWORD"  env-default:"mypassword"`
	Host     string `envconfig:"DB_HOST" env-default:"localhost"`
	Port     string `envconfig:"DB_PORT" env-default:"5430"`
	DBName   string `envconfig:"DB_NAME" env-default:"mydatabase"`
	SSLMode  string `envconfig:"DB_SSLMODE" env-default:"require"`
}
type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Protocol    string        `yaml:"protocol" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func EnvLoad() (*Config, *DataBase) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load environment file, error: ", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	// log.Printf("DataBase configuration: %+v\n", os.Getenv("DB_USERNAME"))
	// os.Setenv("DB_PORT", "5050")

	db, err := EnvLoadDb()
	if err != nil {
		log.Fatal("err db %w\n", err)
	}
	log.Printf("DataBase configuration: %+v\n", db)
	cfg := EnvLoadInPath(configPath)
	return cfg, db
}
func EnvLoadDb() (*DataBase, error) {

	var dbConfig DataBase

	if err := cleanenv.ReadEnv(&dbConfig); err != nil {
		return &dbConfig, fmt.Errorf("failed to load database configuration from environment: %w", err)
	}
	return &dbConfig, nil

}
func EnvLoadInPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config, error: ", err)
	}

	return &cfg
}
