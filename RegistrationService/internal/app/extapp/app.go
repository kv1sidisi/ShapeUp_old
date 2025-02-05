package extapp

import (
	"RegistrationService/cmd/grpccl"
	"RegistrationService/internal/app/intapp"
	"RegistrationService/internal/config"
	"RegistrationService/internal/service/usrcreatesvc"
	"RegistrationService/internal/storage/pgsql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

// App structure represents upper layer of application and configure bottom layer of application with database and register service.
type App struct {
	GRPCSrv *intapp.App
	cfg     *config.Config
}

// New creates upper layer of application
func New(
	log *slog.Logger,
	cfg *config.Config,
	postgresqlClient *pgxpool.Pool,
	grpccl *grpccl.GRPCClients,
) *App {
	storage, err := pgsql.New(postgresqlClient, log)
	if err != nil {
		panic(err)
	}
	log.Info("postgresql storage manager created")

	usrCreateSvc := usrcreatesvc.New(log, storage, grpccl)
	log.Info("user creation service created")

	grpcApp := intapp.New(log, usrCreateSvc, cfg)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
