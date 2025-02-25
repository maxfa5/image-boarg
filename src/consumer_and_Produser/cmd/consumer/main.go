package main

import (
	"fmt"
	consumer "kafka_with_go/internal/Consumer"
	database "kafka_with_go/internal/Dbconnect"
	"kafka_with_go/internal/config"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	file, err := os.Create("loger.txt")
	if err != nil {
		fmt.Println("error in logger")
	}

	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))
	cfg, db_info := config.EnvLoad()
	consumerService, err := consumer.NewConsumerService(logger, cfg.FirstConsumer)
	if err != nil {
		logger.Error("error in create Consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		db_info.Username,
		db_info.Password,
		db_info.Host,
		db_info.Port,
		db_info.DBName,
	)
	log.Println("Connection string:", connectionString)

	err = database.InitDB(logger, connectionString)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	go consumerService.LoopGetMsg()

	cancelCh := make(chan os.Signal, 1)
	signal.Notify(cancelCh, syscall.SIGINT, syscall.SIGTERM)
	stopSignal := <-cancelCh
	consumerService.StopConsumer()
	logger.Info("stoppping server", slog.String("signal", stopSignal.String()))
}
