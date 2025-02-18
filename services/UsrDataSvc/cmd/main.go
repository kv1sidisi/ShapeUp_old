package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	loadconfig "github.com/kv1sidisi/shapeup/pkg/config"
	"github.com/kv1sidisi/shapeup/pkg/database/pgcl"
	"github.com/kv1sidisi/shapeup/pkg/interceptors/authincp"
	"github.com/kv1sidisi/shapeup/pkg/logger"
	pbjwtsvc "github.com/kv1sidisi/shapeup/pkg/proto/jwtsvc/pb"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/cmd/grpccl/consts"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/app/extapp"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := &config.Config{}
	loadconfig.MustLoad(cfg)

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	postgresqlClient := mustConnectToDatabase(cfg)
	log.Info("connected to database")

	//connecting grpc clients
	clients := grpccl.New(log, cfg)
	defer clients.Close()

	//connecting interceptors
	jwtClient := clients.Cl[consts.JWTSvc].Client.(pbjwtsvc.JWTClient)
	authInCp := authincp.AuthInterceptor(log, jwtClient)

	log.Info("interceptors connected")

	application := extapp.New(log, cfg, postgresqlClient, authInCp)

	log.Info("application created")

	go application.GRPCSrv.MustRun()
	log.Info("GRPC server started")

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("received shutdown signal", slog.Any("signal", sign))
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

// mustConnectToDatabase panics if setupDatabaseConnection fails
func mustConnectToDatabase(cfg *config.Config) *pgxpool.Pool {
	postgresqlClient, err := setupDatabaseConnection(cfg)
	if err != nil {
		panic(err)
	}
	return postgresqlClient
}

// setupDatabaseConnection connect to database
// returns pgx client.
func setupDatabaseConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	postgresqlClient, err := pgcl.NewClient(context.TODO(), 3, &cfg.Storage)
	if err != nil {
		return nil, err
	}
	return postgresqlClient, nil
}
