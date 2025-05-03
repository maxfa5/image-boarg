package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
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

func AddMessageToElasticsearch(client *elastic.Client, message MessageData) error {
	ctx := context.Background()

	// Преобразуем сообщение в JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Добавляем сообщение в Elasticsearch
	_, err = client.Index().
		Index("messages").
		Id(fmt.Sprintf("%.0f", message.ChatID)). // Используем ID сообщения как идентификатор документа
		BodyString(string(messageJSON)).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
