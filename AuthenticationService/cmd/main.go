package main

import (
	pbjwtsvc "AuthenticationService/api/pb/jwtsvc"
	pbsendsvc "AuthenticationService/api/pb/sendsvc"
	"AuthenticationService/internal/app/extapp"
	"AuthenticationService/internal/config"
	"AuthenticationService/pkg/client/pgsqlcl"
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

	postgresqlClient := mustConnectToDatabase(cfg, log)
	log.Info("connected to database")

	sendingServiceConn := mustConnectToGRPC(cfg.GRPCClient.SendingServiceAddress, log)
	defer sendingServiceConn.Close()
	sendingClient := pbsendsvc.NewSendingClient(sendingServiceConn)
	log.Info("GRPC SendingService connected", slog.String("address", cfg.GRPCClient.SendingServiceAddress))

	jwtServiceConn := mustConnectToGRPC(cfg.GRPCClient.JWTServiceAddress, log)
	defer jwtServiceConn.Close()
	jwtClient := pbjwtsvc.NewJWTClient(jwtServiceConn)
	log.Info("GRPC JWTService connected", slog.String("address", cfg.GRPCClient.JWTServiceAddress))

	application := extapp.New(log, cfg, postgresqlClient, sendingClient, jwtClient)
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

// mustConnectToDatabase panics if setupDatabaseConnection fails
func mustConnectToDatabase(cfg *config.Config, log *slog.Logger) *pgxpool.Pool {
	postgresqlClient, err := setupDatabaseConnection(cfg, log)
	if err != nil {
		panic(err)
	}
	return postgresqlClient
}

// setupDatabaseConnection connect to database.
func setupDatabaseConnection(cfg *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	postgresqlClient, err := pgsqlcl.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		return nil, err
	}
	return postgresqlClient, nil
}

func mustConnectToGRPC(address string, log *slog.Logger) *grpc.ClientConn {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return conn
}
