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
	storage, err := postgresql.New(postgresqlClient, log)
	if err != nil {
		panic(err)
	}
	log.Info("postgresql storage manager created")

	authService := authsvc.New(log, cfg, storage, sendingClient, jwtClient)
	log.Info("auth service created")

	grpcApp := intapp.New(log, cfg, authService)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
