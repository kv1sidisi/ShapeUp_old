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
	log.Info("creating jwt service")
	jwtService := jwtsvc.New(log, cfg)

	log.Info("creating grpc server")
	grpcApp := intapp.New(log, cfg, jwtService)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
