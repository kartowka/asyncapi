package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	PORT         string `env:"PORT"`
	DB_HOST      string `env:"DB_HOST"`
	DB_PORT      string `env:"DB_PORT"`
	DB_USER      string `env:"DB_USER"`
	DB_PASS      string `env:"DB_PASS"`
	DB_NAME      string `env:"DB_NAME"`
	DATABASE_URL string `env:"DATABASE_URL"`
	JWT_SECRET   string `env:"JWT_SECRET"`
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("%s:%s@/%s?parseTime=true",
		c.DB_USER,
		c.DB_PASS,
		c.DB_NAME,
	)
}
func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
