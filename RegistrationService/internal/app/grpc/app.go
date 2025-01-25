package grpcapp

import (
	"RegistrationService/internal/config"
	reggrpc "RegistrationService/internal/grpc/user_creation"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

// App structure represents bottom layer of application and contains grpc server.
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	cfg        *config.Config
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	registerService reggrpc.UserCreation,
	cfg *config.Config,
) *App {
	gRPCServer := grpc.NewServer()

	// Connects handlers.
	reggrpc.RegisterServer(gRPCServer, registerService, cfg)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		cfg:        cfg,
	}
}

// MustRun runs gRPC server and panics if any errors occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	// Shows where this method is.
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("port", a.cfg.GRPC.Port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	// Start server with listener "l".
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() error {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int64("port", a.cfg.GRPC.Port))

	a.gRPCServer.GracefulStop()

	return nil
}
