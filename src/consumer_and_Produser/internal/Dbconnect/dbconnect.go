package dbconnect

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

// InitDB инициализирует пул подключений к PostgreSQL.
func InitDB(logger *slog.Logger, connectionString string) error {
	var err error
	dbPool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.Error("Не удалось подключиться к базе данных", slog.String("error", err.Error()))
		return fmt.Errorf("Не удалось подключиться к базе данных: %w", err)
	}
	fmt.Println("DB successfully connect")
	return nil
}

// GetDB возвращает пул подключений к базе данных.
func GetDB() *pgxpool.Pool {
	return dbPool
}

// CloseDB закрывает пул подключений.
func CloseDB() {
	dbPool.Close()
}
