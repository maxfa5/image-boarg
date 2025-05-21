package main

import (
	consumer "kafka_with_go/internal/Consumer"
	Elconnect "kafka_with_go/internal/Elasticconnect"
	"kafka_with_go/internal/config"
	"os/signal"
	"syscall"

	"log"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.EnvLoad(logger)
	consumerService, err := consumer.NewConsumerService(logger, cfg.FirstConsumer)
	if err != nil {
		logger.Error("error in create Consumer: ", slog.String("error", err.Error()))

		os.Exit(1)
	}
	err = Elconnect.InitElastic(logger)
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
