package external_app

import (
	internal_app "JWTService/internal/app/grpc"
	"JWTService/internal/config"
	jwt_service "JWTService/internal/service"
	"log/slog"
)

type App struct {
	GRPCSrv *internal_app.App
	cfg     *config.Config
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	log.Info("creating jwt service")
	jwtService := jwt_service.New(log, cfg)

	log.Info("creating grpc server")
	grpcApp := internal_app.New(log, cfg, jwtService)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
