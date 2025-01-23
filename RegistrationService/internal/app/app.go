package app

import (
	grpcapp "RegistrationService/internal/app/grpc"
	"RegistrationService/internal/service/register"
	"RegistrationService/internal/storage/postgresql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int64,
	postgresqlClient *pgxpool.Pool,
) *App {
	storage, err := postgresql.New(postgresqlClient, log)
	if err != nil {
		panic(err)
	}

	registerService := register.New(log, storage)
	grpcApp := grpcapp.New(log, registerService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
