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
	Post_id string  `json: "post_id"`
	Author  string  `json: "author_id"`
	Content string  `json:"content"`
	ChatID  float64 `json:"chat_id"`
	images  []Image `json: "images"`
}

type Image struct {
	URL  string `json:"url"`  // keyword
	Hash string `json:"hash"` // keyword
}

type Images struct {
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

func GetMessagesByChatID(w http.ResponseWriter, r *http.Request, client *elastic.Client) {
	ctx := context.Background()
	vars := mux.Vars(r)
	threadIDStr := vars["thread_id"] // Получаем thread_id из параметров маршрута

	fmt.Printf("vars: %v\n", vars)

	// Преобразование thread_id в целое число
	threadID, err := strconv.ParseInt(threadIDStr, 10, 64) // 10 - base, 64 - bitSize
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный thread_id: %s. Ожидается целое число.", threadIDStr), http.StatusBadRequest)
		return
	}

	fmt.Printf("threadID: %d\n", threadID)

	// Поиск сообщений по thread_id
	query := elastic.NewTermQuery("thread_id", threadID)
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
