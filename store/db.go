package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antfley/asyncapi/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLDb(config *config.Config) (*gorm.DB, error) {
	dsn := config.DatabaseURL()
	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             2 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
