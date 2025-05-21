package main

import (
	"Post_Reader/internal/config"
	"Post_Reader/internal/elastic"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.EnvLoad(logger)
	// Инициализация Elasticsearch клиента

	client, err := elastic.InitElastic(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch client: %v", err)
	}
	// Инициализация роутера
	r := mux.NewRouter()

	getAllMessagesHandler := func(w http.ResponseWriter, r *http.Request) {
		elastic.GetAllMessages(w, r, client)
	}
	getMessagesByChatIDHandler := func(w http.ResponseWriter, r *http.Request) {
		elastic.GetMessagesByChatID(w, r, client)
	}

	getThreads := func(w http.ResponseWriter, r *http.Request) {
		elastic.GetThreads(w, r, client)
	}
	r.HandleFunc("/api/messages", getAllMessagesHandler).Methods("GET")
	r.HandleFunc("/api/threads", getThreads).Methods("GET")
	r.HandleFunc("/api/messages/{thread_id}", getMessagesByChatIDHandler).Methods("GET")
	// Запуск сервера
	fmt.Printf("Сервер запущен на порту %s...\n", cfg.ServerElastic.Port_server)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerElastic.Port_server), r))
}
