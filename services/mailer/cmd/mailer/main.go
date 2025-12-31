package main

import (
	"config"
	"context"
	"eventBus"
	"log"
	"log/slog"
	plog "logger"
	"mailer/internal/handlers"
	"mailer/internal/storage/postgresql"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting mailer service", slog.String("env", cfg.Env))

	storage, err := postgresql.NewStorage("pgx", cfg.ConnectionString)
	if err != nil {
		logger.Error("failed to init storage", plog.Error(err))
		os.Exit(1)
	}

	defer func(storage *sqlx.DB) {
		err = storage.Close()
		logger.Warn("failed to close storage", plog.Error(err))
	}(storage)

	logger.Info("starting application", slog.String("address", cfg.HttpServer.Address))

	consumersContext := context.Background()
	emailsConsumer := eventBus.NewConsumer(cfg.Consumer.Brokers, EmailSendEventTopic, "mailer_service")

	go func() {
		if consumeErr := emailsConsumer.Consume(consumersContext, handlers.NewEmailsHandler(logger).Handle); err != nil {
			logger.Error("consume error", slog.String("error", consumeErr.Error()))
		}
	}()
}

func setupLogger(env string) *slog.Logger {
	logger := &slog.Logger{}

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
