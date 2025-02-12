package main

import (
	loadconfig "github.com/kv1sidisi/shapeup/pkg/config"
	"github.com/kv1sidisi/shapeup/pkg/logger"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/app/extapp"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := &config.Config{}
	loadconfig.MustLoad(cfg)

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

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
