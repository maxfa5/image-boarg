package consumer

import (
	"fmt"
	curd "kafka_with_go/internal/CURD"
	"kafka_with_go/internal/config"
	"log/slog"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerService struct {
	logger   *slog.Logger
	cfg      config.Consumer
	consumer *kafka.Consumer
	stop     chan bool
}

func NewConsumerService(logger *slog.Logger, cfg config.Consumer) (*ConsumerService, error) {

	consumer, err := connToKafkaTopic(cfg)
	if err != nil {
		logger.Error("Ошибка подключения к Kafka", slog.String("error", err.Error()))
	} else {
		logger.Info("Consumer успешно подключен к Kafka" + cfg.Topic + ". Ожидание сообщений...")
	}

	return &ConsumerService{
		logger:   logger,
		cfg:      cfg,
		consumer: consumer,
		stop:     make(chan bool),
	}, err

}

func connToKafkaTopic(cfg config.Consumer) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	for err != nil || consumer == nil {
		consumer, err = kafka.NewConsumer(config)
		if err != nil {
			fmt.Printf("error creating kafka consumer: %v", err)
		}
	}
	_, err = consumer.GetMetadata(nil, true, 500)
	if err != nil {
		consumer.Close()
		fmt.Printf("error retrieving meta: %v", err)
	}

	err = consumer.SubscribeTopics([]string{cfg.Topic}, nil)
	fmt.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAA: %v", cfg.Topic)

	if err != nil {
		consumer.Close()
		fmt.Printf("error subscribing to topic: %v", err)
	}
	return consumer, err
}

func (c *ConsumerService) LoopGetMsg() {
	const op = "consumer.loopGetMsg"
	c.logger = slog.With(
		slog.String("op", op),
	)

	for {
		select {
		case <-c.stop:
			c.logger.Warn("consumer stopped")
			return
		default:
			msg, err := c.consumer.ReadMessage(time.Second * 1)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					c.logger.Debug("No message received yet, continuing loop", slog.String("error", err.Error()))
					continue
				}
				err = fmt.Errorf("error while reading from kafka: %w", err)
				c.logger.Error("Error while reading from kafka", slog.String("error", err.Error()))
				fmt.Printf("error while reading from kafka: %v", err)
				os.Exit(1)
			}
			go curd.HandleKafkaMessage(c.logger, msg)
			c.logger.Info("Message received", slog.String("topic", *msg.TopicPartition.Topic),
				slog.Int("partition", int(msg.TopicPartition.Partition)),
				slog.Any("offset", msg.TopicPartition.Offset),
				slog.String("value", string(msg.Value)))
			// TODO добавть добавление в db
		}

	}
}

func (c *ConsumerService) StopConsumer() {
	c.stop <- true
	c.logger.Info("Consumer is stopped")
}
