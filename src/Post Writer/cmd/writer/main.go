package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
)

type MessageData struct {
	Content string  `json:"content"`
	ChatID  float64 `json:"chat_id"`
}

var client *elastic.Client

func main() {
	// Инициализация Elasticsearch клиента
	var err error
	client, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatalf("Ошибка подключения к Elasticsearch: %v", err)
	}

	// Инициализация роутера
	r := mux.NewRouter()
	r.HandleFunc("/messages", getAllMessages).Methods("GET")
	r.HandleFunc("/messages/{chat_id}", getMessagesByChatID).Methods("GET")

	// Запуск сервера
	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// getAllMessages возвращает все сообщения из Elasticsearch
func getAllMessages(w http.ResponseWriter, r *http.Request) {
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
func getMessagesByChatID(w http.ResponseWriter, r *http.Request) {
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
