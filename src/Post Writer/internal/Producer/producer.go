package producer

import (
	"encoding/json"
	"fmt"
	curd "kafka_with_go/internal/CURD"
	"kafka_with_go/internal/config"
	"log"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerService struct {
	logger   *slog.Logger
	cfg      config.Producer
	producer *kafka.Producer
	stop     chan bool
}

func NewProduserService(logger *slog.Logger, cfg config.Producer) (*ProducerService, error) {

	Producer, err := connToKafkaTopic(cfg)
	if err != nil {
		logger.Error("Ошибка подключения к Kafka", slog.String("error", err.Error()))
	} else {
		logger.Info("Producer успешно подключен к Kafka. Ожидание сообщений...")
	}

	return &ProducerService{
		logger:   logger,
		cfg:      cfg,
		producer: Producer,
		stop:     make(chan bool), //обратить внимание
	}, err

}

func connToKafkaTopic(cfg config.Producer) (*kafka.Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
	}

	producer, err := kafka.NewProducer(config)
	for err != nil || producer == nil {
		producer, err = kafka.NewProducer(config)
		if err != nil {
			fmt.Printf("error creating kafka Producer: %v", err)
		}
	}
	_, err = producer.GetMetadata(nil, true, 500)
	if err != nil {
		producer.Close()
		fmt.Printf("error retrieving meta: %v", err)
	}

	return producer, nil
}

func main() {
	brokers := "localhost:9092,localhost:9093,localhost:9091"

	// Создаем продюсера
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers, // Замените на адрес вашего Kafka брокера
	})
	if err != nil {
		log.Fatalf("Ошибка создания продюсера: %s\n", err)
	}
	defer producer.Close()

	// Отправляем сообщение
	topic := "orders" // Замените на имя вашей темы
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("Привет, Kafka!"),
	}

	// Отправляем сообщение асинхронно
	err = producer.Produce(message, nil)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %s\n", err)
		return
	}

	// Ждем, пока все сообщения будут отправлены
	producer.Flush(15 * 1000)

	fmt.Println("Сообщение успешно отправлено!")
}

func (ps *ProducerService) SendMessageInKafka(logger *slog.Logger, topic string, crudMessage curd.CRUDMessage) error {
	messageValue, err := json.Marshal(crudMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageValue,
	}

	// Отправляем сообщение асинхронно
	err = ps.producer.Produce(message, nil)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %s\n", err)
		return err
	}

	// Ждем, пока все сообщения будут отправлены
	ps.producer.Flush(15 * 1000)

	fmt.Println("Сообщение успешно отправлено в kafka!")
	return nil
}
