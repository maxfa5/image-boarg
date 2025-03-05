package curd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	Elconnect "kafka_with_go/internal/Elasticconnect"
	"log"
	"log/slog"
	"time"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CRUDMessage struct {
	Action string                 `json:"action"`
	Model  string                 `json:"model"`
	Data   map[string]interface{} `json:"data"` // Use a map for flexible data
}

type MessageData struct {
	Content string  `json:"content"`
	ChatID  float64 `json:"chat_id"`
}

func HandleKafkaMessage(logger *slog.Logger, msg *kafka.Message) {
	var crudMessage CRUDMessage
	jsonStr := string(msg.Value)
	err := json.Unmarshal(msg.Value, &crudMessage)
	if err != nil {
		logger.Error("Error unmarshaling JSON", slog.String("json", jsonStr), slog.String("error", err.Error())) // Add JSON to log
		return
	}

	// Determine which model and action to take
	switch crudMessage.Model {
	case "messages":
		handleMessageModel(logger, crudMessage)
	case "users":
		// Handle user-related operations (not implemented here)
		log.Println("Handling users is not yet implemented")
	default:
		log.Printf("Unknown model: %s", crudMessage.Model)
	}
}

func handleMessageModel(logger *slog.Logger, crudMessage CRUDMessage) {
	switch crudMessage.Action {
	case "create":
		var messageData MessageData
		// Convert the data to struct
		jsonStr, _ := json.Marshal(crudMessage.Data)
		json.Unmarshal(jsonStr, &messageData)
		//createMessage(logger, crudMessage.Data, dbconnect.GetDB())
		createMessage(logger, messageData, Elconnect.GetElastic())
	// case "read":
	// 	readMessage(crudMessage.Data)
	// case "update":
	// 	updateMessage(crudMessage.Data)
	// case "delete":
	// 	deleteMessage(crudMessage.Data)
	default:
		log.Printf("Unknown action: %s", crudMessage.Action)
	}
}

func createMessage(logger *slog.Logger, data MessageData, db *elasticsearch.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Преобразуйте структуру MessageData в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal message data to JSON", slog.String("error", err.Error()))
		return
	}
	// Определите индекс, в который вы хотите сохранить документ
	indexName := "messages"

	// Индексируйте документ в Elasticsearch
	resp, err := db.Index(
		indexName,
		bytes.NewReader(jsonData),
		db.Index.WithContext(ctx),
		db.Index.WithRefresh("true"), //  "true" - для немедленного обновления индекса (для целей отладки)
	)
	if err != nil {
		logger.Error("Failed to index document in Elasticsearch", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.IsError() {
		logger.Error("Elasticsearch returned an error", slog.String("status", resp.String()))
		return
	}

	fmt.Println("Successfully indexed message in Elasticsearch")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// // Prepare the statement outside the function (e.g., during application startup)
	// insertMessageSQL := "INSERT INTO messages (content, chat_id) VALUES ($1, $2)"
	// //
	// if _, err := db.Exec(ctx, insertMessageSQL, data.Content, int(data.ChatID)); err != nil {
	// 	logger.Error("Failed to create message", slog.String("error", err.Error()))
	// 	return
	// }
	// fmt.Println("SuccessFully send")
}
