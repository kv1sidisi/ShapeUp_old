package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kv1sidisi/shapeup/services/regsvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/app/extapp"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/regsvc/pkg/client/pgsqlcl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	endProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	postgresqlClient := mustConnectToDatabase(cfg, log)
	log.Info("connected to database")

	//connecting grpc clients
	clients := grpccl.New(log, cfg)
	defer clients.Close()

	application := extapp.New(log, cfg, postgresqlClient, clients)
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

// setupLogger returns slog logger depending on "env".
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case endProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// mustConnectToDatabase panics if setupDatabaseConnection fails
func mustConnectToDatabase(cfg *config.Config, log *slog.Logger) *pgxpool.Pool {
	postgresqlClient, err := setupDatabaseConnection(cfg, log)
	if err != nil {
		panic(err)
	}
	return postgresqlClient
}

// setupDatabaseConnection connect to database
// returns pgx client.
func setupDatabaseConnection(cfg *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	postgresqlClient, err := pgsqlcl.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		return nil, err
	}
	return postgresqlClient, nil
}
