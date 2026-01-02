package main

import (
	"authOrchestrator/internal/orchestrators/registration"
	"authOrchestrator/internal/storage/postgresql"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/flores666/profileshare-lib/config"
	"github.com/flores666/profileshare-lib/eventBus"
	plog "github.com/flores666/profileshare-lib/logger"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")

	if env == "local" || env == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln(err)
		}
	}

	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("starting mailer service", slog.String("env", cfg.Env))

	storage, err := postgresql.NewStorage("pgx", os.Getenv("DB__CONNECTION_STRING"))
	if err != nil {
		logger.Error("failed to init storage", plog.Error(err))
		os.Exit(1)
	}

	defer func(storage *sqlx.DB) {
		err = storage.Close()
		logger.Warn("failed to close storage", plog.Error(err))
	}(storage)

	logger.Info("starting application", slog.String("address", cfg.HttpServer.Address))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	orch := registration.NewOrchestrator(
		eventBus.NewConsumer(cfg.Consumer.Brokers, "users.registered", "authOrchestrator"),
		eventBus.NewProducer(cfg.Producer.Brokers),
		logger,
	)

	go func() {
		if err = orch.Run(ctx); err != nil {
			logger.Error("orchestrator error", plog.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
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
