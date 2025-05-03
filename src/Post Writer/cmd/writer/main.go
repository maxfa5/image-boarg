package main

import (
	"encoding/json"
	"fmt"
	curd "kafka_with_go/internal/CURD"
	database "kafka_with_go/internal/Dbconnect"
	producer "kafka_with_go/internal/Producer"
	"kafka_with_go/internal/config"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	// file, err := os.Create("loger.txt")
	// if err != nil {
	// 	fmt.Println("error in logger")
	// }

	// defer file.Close()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.EnvLoad(logger)
	producerService, err := producer.NewProduserService(logger, cfg.Producer)
	if err != nil {
		logger.Error("error in create Consumer: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = database.InitDB(logger, cfg.DataBase)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	// Инициализация роутера
	r := mux.NewRouter()
	r.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		pushMessage(w, r, producerService, logger)
	}).Methods("POST")

	// Запуск сервера
	fmt.Println("Сервер запущен на порту " + strconv.Itoa(cfg.HTTPServer.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPServer.Port), r))
}

func pushMessage(w http.ResponseWriter, r *http.Request, producerService *producer.ProducerService, logger *slog.Logger) {
	var message curd.CRUDMessage
	fmt.Print(r.Body)
	jsonMsg := json.NewDecoder(r.Body)
	if err := jsonMsg.Decode(&message); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	// log.Printf("%s\n", message.Model)
	curd.HandleMessageInDB(logger, &message)
	go func() {
		err := producerService.SendMessageInKafka(logger, "messages", message)
		if err != nil {
			logger.Error("Ошибка отправки сообщения в Kafka", slog.String("error", err.Error()))
			// http.Error(w, "Ошибка отправки сообщения в Kafka", http.StatusInternalServerError)
			return
		}
	}()
	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Операция выполнена успешно"))
	// if err := curd.HandleKafkaMessage(logger, message); err != nil {
	// 	logger.Error("Ошибка сохранения сообщения в PostgreSQL", slog.String("error", err.Error()))
	// 	http.Error(w, "Ошибка сохранения сообщения", http.StatusInternalServerError)
	// 	return
	// }
}
