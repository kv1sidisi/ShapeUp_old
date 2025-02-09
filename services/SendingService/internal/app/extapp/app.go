package extapp

import (
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/service/sendsvc"
	"log/slog"
)

// App external layer of GRPC application.
type App struct {
	GRPCSrv *intapp.App
	cfg     *config.Config
}

// New creates services and internal GRPC
//
// Returns App.
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
