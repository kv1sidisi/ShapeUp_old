package extapp

import (
	pbjwtsvc "AuthenticationService/api/pb/jwtsvc"
	pbsendsvc "AuthenticationService/api/pb/sendsvc"
	"AuthenticationService/internal/app/intapp"
	"AuthenticationService/internal/config"
	"AuthenticationService/internal/service/authsvc"
	"AuthenticationService/internal/storage/pgsql"
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
	storage, err := postgresql.New(postgresqlClient, log)
	if err != nil {
		log.Error("error creating postgresql service: ", err)
		panic(err)
	}

	log.Info("creating auth service")
	authService := authsvc.New(log, cfg, storage)

	log.Info("creating grpc server external_app")
	grpcApp := intapp.New(log, cfg, authService, sendingClient, jwtClient)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
