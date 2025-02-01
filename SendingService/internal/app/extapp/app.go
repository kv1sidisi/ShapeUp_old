package extapp

import (
	"SendingService/internal/app/intapp"
	"SendingService/internal/config"
	"SendingService/internal/service/sendsvc"
	"log/slog"
)

// App structure represents upper layer of application and configure bottom layer of application with database and register service.
type App struct {
	GRPCSrv *intapp.App
	cfg     *config.Config
}

// New creates upper layer of application
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	log.Info("creating email sending service")
	sendingService := sendsvc.New(log, cfg)

	log.Info("creating grpc server app")
	grpcApp := intapp.New(log, sendingService, cfg)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
