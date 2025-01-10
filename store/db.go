package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/antfley/asyncapi/config"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func NewMySQLDb(config *config.Config) (*sql.DB, error) {
	dsn := config.DatabaseURL()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
