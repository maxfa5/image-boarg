package Elconnect

import (
	"fmt"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
)

var esClient *elasticsearch.Client

// InitElastic инициализирует клиент Elasticsearch.
func InitElastic(logger *slog.Logger) error {
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%s", "elasticsearch", "9200")}, // or "https://..."
		// Username:  es_info.Username,                                                  // Optional
		// Password:  es_info.Password,                                                  // Optional
		// ... Other configuration optionsls
	}

	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error("Error creating the Elasticsearch client", slog.String("error", err.Error()))
		return fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error("Error creating the Elasticsearch client", slog.String("error", err.Error()))
		return fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	// Test the connection
	res, err := esClient.Info()
	if err != nil {
		logger.Error("Error getting Elasticsearch info", slog.String("error", err.Error()))
		return fmt.Errorf("error getting Elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("Error response from Elasticsearch", slog.String("status", res.String()))
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	logger.Info("Successfully connected to Elasticsearch", slog.String("version", res.String()))
	fmt.Println("Elasticsearch successfully connected")

	return nil
}

// GetElastic возвращает клиент Elasticsearch.
func GetElastic() *elasticsearch.Client {
	return esClient
}
