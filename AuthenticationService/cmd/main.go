package main

import (
	pbjwt "AuthenticationService/api/pb/jwt_service"
	pbsendsvc "AuthenticationService/api/pb/sending_service"
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

const (
	sendingService = "sending service"
	jwtService     = "jwt service"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	log.Info("connecting to database")
	postgresqlClient := mustConnectToDatabase(cfg, log)
	log.Info("connected to postgresql")

	log.Info("connecting to grpc SendingService", slog.String("address", cfg.GRPCClient.SendingServiceAddress))
	sendingServiceConn := mustConnectToGRPC(cfg.GRPCClient.SendingServiceAddress, log, sendingService)
	defer sendingServiceConn.Close()
	log.Info("grpc sendingClient connected")
	sendingClient := pbsendsvc.NewSendingClient(sendingServiceConn)

	log.Info("connecting to grpc JWTService", slog.String("address", cfg.GRPCClient.JWTServiceAddress))
	jwtServiceConn := mustConnectToGRPC(cfg.GRPCClient.JWTServiceAddress, log, jwtService)
	defer jwtServiceConn.Close()
	log.Info("grpc sendingClient connected")
	jwtClient := pbjwt.NewJWTClient(jwtServiceConn)

	application := extapp.New(log, cfg, postgresqlClient, sendingClient, jwtClient)

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
		log.Error("Failed to connect to database", err)
		return nil, err
	}
	return postgresqlClient, nil
}

func mustConnectToGRPC(address string, log *slog.Logger, serviceName string) *grpc.ClientConn {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create grpc connection to "+serviceName+" service: ", err)
		panic(err)
	}
	return conn
}
