package extapp

import (
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/service/jwtsvc"
	"log/slog"
)

type App struct {
	GRPCSrv *intapp.App
	cfg     *config.Config
}

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
