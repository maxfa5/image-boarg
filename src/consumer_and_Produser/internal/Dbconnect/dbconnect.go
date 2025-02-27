package dbconnect

import (
	"context"
	"fmt"
	"kafka_with_go/internal/config"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

// InitDB инициализирует пул подключений к PostgreSQL.
func InitDB(logger *slog.Logger, db_info config.DataBase) error {
	var err error
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		db_info.Username,
		db_info.Password,
		db_info.Host,
		db_info.Port,
		db_info.DBName,
	)
	log.Println("Connection string:", connectionString)

	dbPool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.Error("Не удалось подключиться к базе данных", slog.String("error", err.Error()))
		return err
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
