package main

import (
	"GatewayAPI/internal/config"
	"GatewayAPI/internal/http-server/middleware/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("%+v\n", cfg)

	log := setupLogger(cfg.Env)

	router := chi.NewRouter()

	// Adds id to requests
	router.Use(middleware.RequestID)
	// Logs requests
	router.Use(mwLogger.New(log))
	// Recovers from panic of middleware
	router.Use(middleware.Recoverer)
	// Middleware to write pretty URLs
	router.Use(middleware.URLFormat)

}

// setupLogger initializes logger dependent on environment
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
