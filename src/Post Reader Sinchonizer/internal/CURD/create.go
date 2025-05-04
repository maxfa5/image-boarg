package curd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	Elconnect "kafka_with_go/internal/Elasticconnect"
	"log"
	"log/slog"
	"strings"
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
	PostID    string    `json:"post_id"`   // keyword
	ThreadID  string    `json:"thread_id"` // keyword
	AuthorID  string    `json:"author_id"` // keyword (ссылка на PostgreSQL users.id/user.name)
	Content   string    `json:"content"`   // text
	Images    []Image   `json:"images"`    // nested
	Timestamp time.Time `json:"timestamp"` // date
}

type Image struct {
	URL  string `json:"url"`  // keyword
	Hash string `json:"hash"` // keyword
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

// func createMessage(logger *slog.Logger, data MessageData, db *elasticsearch.Client) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Преобразуйте структуру MessageData в JSON
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		logger.Error("Failed to marshal message data to JSON", slog.String("error", err.Error()))
// 		return
// 	}
// 	// Определите индекс, в который вы хотите сохранить документ
// 	indexName := "messages"

// 	// Индексируйте документ в Elasticsearch
// 	resp, err := db.Index(
// 		indexName,
// 		bytes.NewReader(jsonData),
// 		db.Index.WithContext(ctx),
// 		db.Index.WithRefresh("true"), //  "true" - для немедленного обновления индекса (для целей отладки)
// 	)
// 	if err != nil {
// 		logger.Error("Failed to index document in Elasticsearch", slog.String("error", err.Error()))
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.IsError() {
// 		logger.Error("Elasticsearch returned an error", slog.String("status", resp.String()))
// 		return
// 	}

// 	fmt.Println("Successfully indexed message in Elasticsearch")

// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()
// 	// // Prepare the statement outside the function (e.g., during application startup)
// 	// insertMessageSQL := "INSERT INTO messages (content, chat_id) VALUES ($1, $2)"
// 	// //
// 	// if _, err := db.Exec(ctx, insertMessageSQL, data.Content, int(data.ChatID)); err != nil {
// 	// 	logger.Error("Failed to create message", slog.String("error", err.Error()))
// 	// 	return
// 	// }
// 	// fmt.Println("SuccessFully send")
// }

func createMessage(logger *slog.Logger, data MessageData, client *elasticsearch.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Устанавливаем timestamp если не задан
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal message data",
			slog.String("error", err.Error()))
		return
	}

	resp, err := client.Index(
		"messages",
		bytes.NewReader(jsonData),
		client.Index.WithContext(ctx),
		client.Index.WithRefresh("true"),
		client.Index.WithDocumentID(data.PostID), // Используем post_id как ID документа
	)
	if err != nil {
		logger.Error("Failed to index message",
			slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.IsError() {
		logger.Error("Elasticsearch error",
			slog.String("status", resp.String()))
		return
	}

	logger.Info("Message indexed successfully",
		slog.String("post_id", data.PostID),
		slog.String("thread_id", data.ThreadID))
}

func CreateMessagesIndex(client *elasticsearch.Client) error {
	mapping := `{
        "mappings": {
            "properties": {
                "post_id":    { "type": "keyword" },
                "thread_id":  { "type": "keyword" },
                "author_id":  { "type": "keyword" },
                "content":    { 
                    "type": "text",
                    "analyzer": "russian" 
                },
                "images": {
                    "type": "nested",
                    "properties": {
                        "url":  { "type": "keyword" },
                        "hash": { "type": "keyword" }
                    }
                },
                "timestamp": { "type": "date" }
            }
        },
        "settings": {
            "analysis": {
                "analyzer": {
                    "russian": {
                        "type": "custom",
                        "tokenizer": "standard",
                        "filter": ["lowercase", "russian_morphology"]
                    }
                }
            }
        }
    }`

	resp, err := client.Indices.Create(
		"messages",
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("error response: %s", resp.String())
	}

	return nil
}
