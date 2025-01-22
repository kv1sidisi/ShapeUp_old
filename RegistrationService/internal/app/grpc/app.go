package grpc

import (
	reggrpc "RegistrationService/internal/grpc/register"
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int64
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	port int64,
) *App {
	gRPCServer := grpc.NewServer()

	// Connects handlers
	reggrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
