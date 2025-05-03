package curd

import (
	"context"
	"encoding/json"
	"fmt"
	dbconnect "kafka_with_go/internal/Dbconnect"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func HandleMessageInDB(logger *slog.Logger, msg *CRUDMessage) {
	// var crudMessage CRUDMessage
	// jsonStr := string(msg.Value)
	// err := json.Unmarshal(msg.Value, &crudMessage)
	// if err != nil {
	// 	logger.Error("Error unmarshaling JSON", slog.String("json", jsonStr), slog.String("error", err.Error())) // Add JSON to log
	// 	return
	// }
	if len(msg.Data) == 0 {
		logger.Error("Error empty value in message")
		return
	}

	// Determine which model and action to take
	switch msg.Model {
	case "messages":
		handleMessageModel(logger, *msg)
	case "users":
		// Handle user-related operations (not implemented here)
		log.Println("Handling users is not yet implemented")
	default:
		log.Printf("Unknown model: %s", msg.Model)
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
		createMessage(logger, messageData, dbconnect.GetDB())
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

func createMessage(logger *slog.Logger, data MessageData, db *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	const insertMessageSQL = "INSERT INTO messages (content, chat_id) VALUES ($1, $2)"
	//
	if _, err := db.Exec(ctx, insertMessageSQL, data.Content, int(data.ChatID)); err != nil {
		logger.Error("Failed to create message", slog.String("error", err.Error()))
		return
	}

	fmt.Println("Successfully create message in PostgeSQL")
}
