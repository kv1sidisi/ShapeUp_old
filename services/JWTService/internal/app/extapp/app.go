package extapp

import (
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/service/jwtsvc"
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
	jwtService := jwtsvc.New(log, cfg)
	log.Info("jwt service created")

	grpcApp := intapp.New(log, cfg, jwtService)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
