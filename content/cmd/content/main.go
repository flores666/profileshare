package main

import (
	"content/internal/config"
	"content/internal/handlers/content"
	customMiddleware "content/internal/lib/logger/middleware"
	"content/internal/lib/logger/sl"
	"content/internal/storage/postgresql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	logger.Info("starting content service", slog.String("env", cfg.Env))

	storage, err := postgresql.NewStorage("pgx", cfg.ConnectionString)
	if err != nil {
		logger.Error("failed to init storage", sl.Error(err))
		os.Exit(1)
	}

	defer func(storage *sqlx.DB) {
		err := storage.Close()
		logger.Warn("failed to close storage", sl.Error(err))
	}(storage)

	server := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      buildHandler(logger, storage),
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IddleTimeout,
	}

	logger.Info("starting application", slog.String("address", cfg.HttpServer.Address))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start http server", sl.Error(err))
	}

	logger.Info("http server stopped")
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

func buildHandler(logger *slog.Logger, storage *sqlx.DB) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(customMiddleware.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	content.NewContentHandler(
		content.NewService(content.NewRepository(storage), logger),
		logger,
	).RegisterRoutes(router)

	return router
}
