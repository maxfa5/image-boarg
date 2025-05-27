package elastic

import (
	"Post_Reader/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
)

type MessageData struct {
	Post_id      string    `json:"post_id"`
	Author       string    `json:"author_id"`
	Content      string    `json:"content"`
	ChatID       float64   `json:"chat_id"`
	Images       []Image   `json:"images"`
	Timestamp    time.Time `json:"timestamp"`      // date
	IsThreadRoot bool      `json:"is_thread_root"` // boolean (true если сообщение начало треда)
}
type ThreadData struct {
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
	fmt.Println("aaaa" + r.Method)
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
func GetThreads(w http.ResponseWriter, r *http.Request, client *elastic.Client) {
	ctx := context.Background()

	// Добавляем сортировку по времени создания (новые треды первыми)
	searchResult, err := client.Search().
		Index("threads").
		Sort("created_at", false). // Сортировка по убыванию (новые первыми)
		Size(100).                 // Ограничиваем количество результатов
		Do(ctx)

	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при поиске тредов: %v", err), http.StatusInternalServerError)
		return
	}

	// Проверяем, есть ли результаты
	if searchResult.Hits == nil || len(searchResult.Hits.Hits) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ThreadData{}) // Возвращаем пустой массив
		return
	}

	// Преобразование результатов
	threads := make([]ThreadData, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var thread ThreadData
		if err := json.Unmarshal(hit.Source, &thread); err != nil {
			log.Printf("Ошибка декодирования треда ID %s: %v", hit.Id, err)
			continue // Пропускаем битые записи вместо прерывания
		}
		threads = append(threads, thread)
	}

	// Отправка результата
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(threads); err != nil {
		log.Printf("Ошибка кодирования ответа: %v", err)
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func GetMessagesByChatID(w http.ResponseWriter, r *http.Request, client *elastic.Client) {
	ctx := context.Background()
	vars := mux.Vars(r)
	threadIDStr := vars["thread_id"]

	fmt.Printf("Searching messages for thread_id: %s\n", threadIDStr)

	// Поиск сообщений по thread_id
	query := elastic.NewTermQuery("thread_id", threadIDStr)
	searchResult, err := client.Search().
		Index("messages").
		Query(query).
		Size(100).
		Do(ctx)

	if err != nil {
		http.Error(w, "Ошибка при поиске сообщений: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразование результатов в структуру MessageData
	var messages []MessageData
	for _, hit := range searchResult.Hits.Hits {
		var msg MessageData
		if hit.Source != nil {
			err := json.Unmarshal(hit.Source, &msg)
			if err != nil {
				http.Error(w, "Ошибка при декодировании сообщения: "+err.Error(), http.StatusInternalServerError)
				return
			}
			messages = append(messages, msg)
		}
	}

	// Логирование количества найденных сообщений
	fmt.Printf("Found %d messages for thread_id %s\n", len(messages), threadIDStr)

	// Отправка результата один раз после цикла
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Ошибка при кодировании ответа: "+err.Error(), http.StatusInternalServerError)
	}
}
