package app

import (
	"RegistrationService/api/pb/sending_service"
	grpcapp "RegistrationService/internal/app/grpc"
	"RegistrationService/internal/config"
	"RegistrationService/internal/service/user_creation"
	"RegistrationService/internal/storage/postgresql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

// App structure represents upper layer of application and configure bottom layer of application with database and register service.
type App struct {
	GRPCSrv *grpcapp.App
	cfg     *config.Config
}

// New creates upper layer of application
func New(
	log *slog.Logger,
	cfg *config.Config,
	postgresqlClient *pgxpool.Pool,
	sendingClient sending_service.SendingClient,
) *App {
	log.Info("creating postgresql service")
	storage, err := postgresql.New(postgresqlClient, log)
	if err != nil {
		log.Error("error creating postgresql service: ", err)
		panic(err)
	}

	log.Info("creating register service")
	registerService := user_creation.New(log, storage)

	log.Info("creating grpc server app")
	grpcApp := grpcapp.New(log, registerService, cfg, sendingClient)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
