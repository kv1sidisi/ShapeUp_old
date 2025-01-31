package internal_app

import (
	"JWTService/internal/config"
	grpc_server "JWTService/internal/grpc"
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	log        *slog.Logger
	cfg        *config.Config
	grpcServer *grpc.Server
}

func New(log *slog.Logger,
	cfg *config.Config,
	jwtService grpc_server.JWT,
) *App {
	gRPCServer := grpc.NewServer()
	log.Info("grpc server created")

	log.Info("registering services in grpc server")
	grpc_server.RegisterServer(gRPCServer, jwtService, cfg, log)

	return &App{log: log, cfg: cfg}
}
