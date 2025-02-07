package extapp

import (
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/service/sendsvc"
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
	sendingService := sendsvc.New(log, cfg)
	log.Info("email sending service created")

	grpcApp := intapp.New(log, sendingService, cfg)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
