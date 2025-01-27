package main

import (
	pb "RegistrationService/api/pb/sending_service"
	"RegistrationService/internal/app"
	"RegistrationService/internal/config"
	"RegistrationService/pkg/client/postgresql"
	"context"
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

	log.Info("connecting to database")
	postgresqlClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		log.Error("Failed to connect to database", err)
		panic(err)
	}
	log.Info("connected to database")

	log.Info("connecting to grpc SendingService", slog.String("address", cfg.GRPCClient.SendingServiceAddress))
	conn, err := grpc.NewClient(cfg.GRPCClient.SendingServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create grpc client connection to sending service: ", err)
		panic(err)
	}
	defer conn.Close()
	log.Info("grpc client connected")

	client := pb.NewSendingClient(conn)

	log.Info("starting application")
	application := app.New(log, cfg, postgresqlClient, client)

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
