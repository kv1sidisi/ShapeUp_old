package app

import (
	grpcapp "SendingService/internal/app/grpc"
	"SendingService/internal/config"
	"SendingService/internal/service"
	"log/slog"
)

// App structure represents upper layer of application and configure bottom layer of application with database and register service.
type App struct {
	GRPCSrv *grpcapp.App
	cfg     *config.Config
}

// New creates upper layer of application
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	log.Info("creating email sending service")
	sendingService := service.New(log, cfg)

	log.Info("creating grpc server app")
	grpcApp := grpcapp.New(log, sendingService, cfg)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
