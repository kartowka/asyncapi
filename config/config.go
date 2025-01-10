package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	DB_HOST      string
	DB_PORT      string
	DB_USER      string
	DB_PASS      string
	DB_NAME      string
	DATABASE_URL string
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("%s:%s@/%s?parseTime=true",
		c.DB_USER,
		c.DB_PASS,
		c.DB_NAME,
	)
}
func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Config{
		PORT:    os.Getenv("PORT"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_NAME: os.Getenv("DB_NAME"),
	}, nil
}
