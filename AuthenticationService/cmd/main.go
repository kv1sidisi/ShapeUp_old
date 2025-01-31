package main

import (
	pb "AuthenticationService/api/pb/sending"
	external_app "AuthenticationService/internal/app"
	"AuthenticationService/internal/config"
	"AuthenticationService/pkg/client/postgresql"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	log.Info("connecting to database")
	postgresqlClient := mustLoadDatabaseConnection(cfg, log)
	log.Info("connected to postgresql")

	log.Info("connecting to grpc SendingService", slog.String("address", cfg.GRPCClient.SendingServiceAddress))
	conn, err := grpc.NewClient(cfg.GRPCClient.SendingServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create grpc client connection to sending service: ", err)
		panic(err)
	}
	defer conn.Close()
	log.Info("grpc client connected")
	sendingClient := pb.NewSendingClient(conn)

	application := external_app.New(log, cfg, postgresqlClient, sendingClient)

	log.Info("starting grpc server")
	go application.GRPCSrv.MustRun()
	log.Info("grpc server started")

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("received shutdown signal", slog.Any("signal", sign))
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// mustLoadDatabaseConnection panics if setupDatabaseConnection fails
func mustLoadDatabaseConnection(cfg *config.Config, log *slog.Logger) *pgxpool.Pool {
	postgresqlClient, err := setupDatabaseConnection(cfg, log)
	if err != nil {
		panic(err)
	}
	return postgresqlClient
}

// setupDatabaseConnection connect to database.
func setupDatabaseConnection(cfg *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	postgresqlClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		log.Error("Failed to connect to database", err)
		return nil, err
	}
	return postgresqlClient, nil
}
