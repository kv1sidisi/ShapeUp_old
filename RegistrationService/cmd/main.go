package main

import (
	pbjwtsvc "RegistrationService/api/grpccl/pb/jwtsvc"
	pbsendsvc "RegistrationService/api/grpccl/pb/sendsvc"
	"RegistrationService/internal/app/extapp"
	"RegistrationService/internal/config"
	"RegistrationService/pkg/client/pgsqlcl"
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

// setupLogger sets up logger dependent on environment type.
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
