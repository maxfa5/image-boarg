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
	cfg := config.EnvLoad(logger)
	consumerService, err := consumer.NewConsumerService(logger, cfg.FirstConsumer)
	if err != nil {
		logger.Error("error in create Consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = database.InitDB(logger, cfg.DataBase)
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
