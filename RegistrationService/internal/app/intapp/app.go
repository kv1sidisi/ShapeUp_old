package intapp

import (
	pbjwtsvc "RegistrationService/api/pb/jwtsvc"
	pbsendsvc "RegistrationService/api/pb/sendsvc"
	"RegistrationService/internal/config"
	"RegistrationService/internal/grpc/grpcsrv"
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

// New creates new gRPC server external_app.
func New(
	log *slog.Logger,
	registerService grpcsrv.UsrCreationSvc,
	cfg *config.Config,
	sendingClient pbsendsvc.SendingClient,
	jwtClient pbjwtsvc.JWTClient,
) *App {
	gRPCServer := grpc.NewServer()
	log.Info("grpc server created")

	log.Info("registering services in grpc server")
	grpcsrv.RegisterServer(gRPCServer, registerService, cfg, log, sendingClient, jwtClient)

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
		log.Error("failed to listen", err)
		return err
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	// Start server with listener "l".
	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("failed to serve", err)
		return err
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
