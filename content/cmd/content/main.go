package main

import (
	"content/internal/handlers/content"
	"content/internal/storage/postgresql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/flores666/profileshare-lib/config"
	libmiddleware "github.com/flores666/profileshare-lib/middleware"

	plog "github.com/flores666/profileshare-lib/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	logger.Info("starting content service", slog.String("env", cfg.Env))

	storage, err := postgresql.NewStorage("pgx", os.Getenv("DB__CONNECTION_STRING"))
	if err != nil {
		logger.Error("failed to init storage", plog.Error(err))
		os.Exit(1)
	}

	defer func(storage *sqlx.DB) {
		err = storage.Close()
		if err != nil {
			logger.Warn("failed to close storage", plog.Error(err))
		}
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
		logger.Error("failed to start http server", plog.Error(err))
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
	router.Use(plog.NewRequestLogMiddleware(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	libmiddleware.AuthMiddleware([]byte(os.Getenv("SECURITY__ACCESS_SECRET")))

	content.NewContentHandler(content.NewService(content.NewRepository(storage), logger)).RegisterRoutes(router)

	return router
}
