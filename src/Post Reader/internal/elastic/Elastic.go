package elastic

import (
	"Post_Reader/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
)

type MessageData struct {
	Content string  `json:"content"`
	ChatID  float64 `json:"chat_id"`
}

var client *elastic.Client

func InitElastic(cfg config.Config) (*elastic.Client, error) {

	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", cfg.ServerElastic.Host, cfg.ServerElastic.Port),
	}
	var err error
	client, err = elastic.NewClient(elastic.SetURL(u.String()))
	if err != nil {
		log.Printf("Ошибка подключения к Elasticsearch: %v", err)             // Log the error
		return nil, fmt.Errorf("ошибка подключения к Elasticsearch: %w", err) // Wrap for context
	}

	return client, nil
}

// getAllMessages возвращает все сообщения из Elasticsearch
func GetAllMessages(w http.ResponseWriter, r *http.Request, client *elastic.Client) {
	ctx := context.Background()

	// Поиск всех сообщений
	searchResult, err := client.Search().
		Index("messages"). // Укажите имя вашего индекса
		Do(ctx)
	if err != nil {
		http.Error(w, "Ошибка при поиске сообщений", http.StatusInternalServerError)
		return
	}

	// Преобразование результатов в структуру MessageData
	var messages []MessageData
	for _, hit := range searchResult.Hits.Hits {
		var msg MessageData
		err := json.Unmarshal(hit.Source, &msg)
		if err != nil {
			http.Error(w, "Ошибка при декодировании сообщения", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	// Отправка результата в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// getMessagesByChatID возвращает сообщения по chat_id
func GetMessagesByChatID(w http.ResponseWriter, r *http.Request, client *elastic.Client) {
	ctx := context.Background()
	vars := mux.Vars(r)
	chatIDStr := vars["chat_id"]

	// Преобразование chat_id в число
	chatID, err := strconv.ParseFloat(chatIDStr, 64)
	if err != nil {
		http.Error(w, "Некорректный chat_id", http.StatusBadRequest)
		return
	}

	// Поиск сообщений по chat_id
	query := elastic.NewTermQuery("chat_id", chatID)
	searchResult, err := client.Search().
		Index("messages"). // Укажите имя вашего индекса
		Query(query).
		Size(100).
		Do(ctx)
	if err != nil {
		http.Error(w, "Ошибка при поиске сообщений", http.StatusInternalServerError)
		return
	}

	// Преобразование результатов в структуру MessageData
	var messages []MessageData
	for _, hit := range searchResult.Hits.Hits {
		var msg MessageData
		err := json.Unmarshal(hit.Source, &msg)
		if err != nil {
			http.Error(w, "Ошибка при декодировании сообщения", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	// Отправка результата в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
