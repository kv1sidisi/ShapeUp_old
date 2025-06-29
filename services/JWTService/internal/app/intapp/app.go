package intapp

import (
	"fmt"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/grpc/grpcsrv"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

// App internal layer of GRPC application.
type App struct {
	log        *slog.Logger
	cfg        *config.Config
	grpcServer *grpc.Server
}

// New creates GRPC server and registers services.
//
// Returns App
func New(log *slog.Logger,
	cfg *config.Config,
	jwtService grpcsrv.JWTSvc,
) *App {
	gRPCServer := grpc.NewServer()
	log.Info("GRPC server created")

	grpcsrv.RegisterServer(gRPCServer, jwtService, cfg, log)
	log.Info("services registered in GRPC server")

	return &App{log: log,
		cfg:        cfg,
		grpcServer: gRPCServer}
}

// MustRun tries to run GRPC server.
//
// Panics if any errors occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs GRPC server.
//
// Returns:
//   - Error if: Fails to listen TCP port. Error while serving GRPC requests.
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

	log.Info("GRPC server is running", slog.String("addr", l.Addr().String()))

	// Start server with listener "l".
	if err := a.grpcServer.Serve(l); err != nil {
		log.Error("failed to serve", err)
		return err
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() error {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping GRPC server", slog.Int64("port", a.cfg.GRPC.Port))

	a.grpcServer.GracefulStop()

	return nil
}
