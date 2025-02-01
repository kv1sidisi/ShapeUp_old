package extapp

import (
	pbjwtsvc "RegistrationService/api/pb/jwtsvc"
	pbsendsvc "RegistrationService/api/pb/sendsvc"
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
	sendingClient pbsendsvc.SendingClient,
	jwtClient pbjwtsvc.JWTClient,
) *App {
	log.Info("creating postgresql service")
	storage, err := pgsql.New(postgresqlClient, log)
	if err != nil {
		log.Error("error creating postgresql service: ", err)
		panic(err)
	}

	log.Info("creating register service")
	registerService := usrcreatesvc.New(log, storage)

	log.Info("creating grpc server external_app")
	grpcApp := intapp.New(log, registerService, cfg, sendingClient, jwtClient)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
