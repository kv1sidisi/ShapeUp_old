package app

import (
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
) *App {
	storage, err := postgresql.New(postgresqlClient, log)
	if err != nil {
		panic(err)
	}

	registerService := user_creation.New(log, storage)
	grpcApp := grpcapp.New(log, registerService, cfg)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
