package main

import (
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/app/extapp"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := extapp.New(log, cfg)
	log.Info("application created")

	go application.GRPCSrv.MustRun()
	log.Info("GRPC server started")

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("received shutdown signal", slog.Any("signal", sign))
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case prod:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
