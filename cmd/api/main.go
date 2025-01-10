package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/antfley/asyncapi/api/server"
	"github.com/antfley/asyncapi/config"
	"github.com/antfley/asyncapi/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg, err := config.New()
	if err != nil {
		return err
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	db, err := store.NewDb(cfg)
	if err != nil {
		return err
	}
	dataStore := store.New(db)
	server := server.New(cfg, logger, dataStore)
	if err := server.Run(ctx); err != nil {
		return err
	}
	return nil
}
