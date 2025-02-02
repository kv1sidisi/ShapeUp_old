package extapp

import (
	"JWTService/internal/app/intapp"
	"JWTService/internal/config"
	"JWTService/internal/service/jwtsvc"
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
