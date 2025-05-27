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
	"github.com/google/uuid"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CRUDMessage struct {
	Action string                 `json:"action"`
	Model  string                 `json:"model"`
	Data   map[string]interface{} `json:"data"` // Use a map for flexible data
}

type MessageData struct {
	PostID       string    `json:"post_id"`        // keyword
	ThreadID     string    `json:"thread_id"`      // keyword
	AuthorID     string    `json:"author_id"`      // keyword (ссылка на PostgreSQL users.id/user.name)
	Content      string    `json:"content"`        // text
	Images       []Image   `json:"images"`         // nested
	Timestamp    time.Time `json:"timestamp"`      // date
	IsThreadRoot bool      `json:"is_thread_root"` // boolean (true если сообщение начало треда)
}
type ThreadDocument struct {
	ThreadID   string    `json:"thread_id"`
	Title      string    `json:"title"`
	RootPostID string    `json:"root_post_id"`
	IsClosed   bool      `json:"is_closed"`
	CreatedAt  time.Time `json:"created_at"`
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

func createMessage(logger *slog.Logger, data MessageData, client *elasticsearch.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Устанавливаем timestamp если не задан
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now().UTC()
	}

	if data.PostID == "" {
		newUUID := uuid.New()
		data.PostID = newUUID.String()
	}
	fmt.Println(data.Timestamp)

	if data.IsThreadRoot {
		var err error
		data.ThreadID, err = createNewThread(logger, data.PostID, data.Content, client)
		if err != nil {
			logger.Error("Failed to index message",
				slog.String("error", err.Error()))
			return
		}

	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal message data",
			slog.String("error", err.Error()))
		return
	}

	// Создаем индекс, если он не существует
	if err := CreateIndexIfNotExists(client, "messages"); err != nil {
		logger.Error("Elasticsearch error in create index")
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

func createNewThread(logger *slog.Logger, postID string, threadTitle string, client *elasticsearch.Client) (string, error) {
	threadID := uuid.New().String()

	// Создаем документ треда
	thread := ThreadDocument{
		ThreadID:   threadID,
		RootPostID: postID,
		Title:      threadTitle,
		IsClosed:   false,
		CreatedAt:  time.Now(),
	}

	// Преобразуем структуру в JSON
	doc, err := json.Marshal(thread)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	// Индексируем документ
	resp, err := client.Index(
		"threads",
		bytes.NewReader(doc),
		client.Index.WithDocumentID(threadID),
		client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return "", fmt.Errorf("elasticsearch error: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		logger.Error("Elasticsearch error",
			slog.String("status", resp.String()))
		return "", fmt.Errorf("Elasticsearch error")
	}

	logger.Info("New thread created",
		slog.String("thread_id", threadID),
		slog.String("root_post_id", postID))

	return threadID, nil
}

func CreateIndexIfNotExists(client *elasticsearch.Client, indexName string) error {
	// Проверяем, существует ли индекс
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		// Индекс не существует, создаем его
		mapping := `{
		"mappings": {
			"properties": {
				"author_id": {
					"type": "keyword"
				},
				"chat_id": {
					"type": "long"
				},
				"content": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"images": {
					"type": "nested",
					"properties": {
						"hash": {
							"type": "text",
							"fields": {
								"keyword": {
									"type": "keyword",
									"ignore_above": 256
								}
							}
						},
						"url": {
							"type": "text",
							"fields": {
								"keyword": {
									"type": "keyword",
									"ignore_above": 256
								}
							}
						}
					}
				},
				"is_thread_root": {
					"type": "boolean"
				},
				"post_id": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"thread_id": {
					"type": "keyword"
				},
				"timestamp": {
					"type": "date"
				}
			}
		}
	}`
		// Создаем индекс с указанным маппингом
		createIndexResponse, err := client.Indices.Create(
			"messages",
			client.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		defer createIndexResponse.Body.Close()

		if createIndexResponse.IsError() {
			return fmt.Errorf("error creating index: %s", createIndexResponse.String())
		}
		log.Printf("Index %s created successfully", indexName)
	} else {
		log.Printf("Index %s already exists", indexName)
	}

	return nil
}
