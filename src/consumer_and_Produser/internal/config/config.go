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
	FirstConsumer Consumer
	HTTPServer
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
type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Protocol    string        `yaml:"protocol" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func EnvLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load environment file, error: ", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	cfg := EnvLoadInPath(configPath)
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
func EnvLoadInPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg.HTTPServer); err != nil { // Загружаем только HTTPServer
		fmt.Printf("failed to read config, error: %+v\n", err) // Возвращаем ошибку, а не вызываем log.Fatal
	}
	if err := cleanenv.ReadEnv(&cfg.FirstConsumer); err != nil {
		log.Fatalf("failed to load consumer configuration from environment: %v\n", err)
	}

	var db *DataBase
	db, err := EnvLoadDb()
	if err != nil {
		log.Fatal("err db %w\n", err)
	}
	cfg.DataBase = *db
	log.Printf("DataBase configuration: %+v\n", cfg.DataBase)
	return &cfg
}
