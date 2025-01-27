package app

import (
	grpcapp "SendingService/internal/app/grpc"
	"SendingService/internal/config"
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
	//TODO: create email sending service

	log.Info("creating grpc server app")
	grpcApp := grpcapp.New(log, registerService, cfg)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
