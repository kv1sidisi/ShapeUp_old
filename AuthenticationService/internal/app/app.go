package external_app

import (
	"AuthenticationService/api/pb/jwt_service"
	"AuthenticationService/api/pb/sending_service"
	internal "AuthenticationService/internal/app/grpc"
	"AuthenticationService/internal/config"
	"AuthenticationService/internal/service/auth_service"
	"AuthenticationService/internal/storage/postgresql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

// App structure represents upper layer of application and configure bottom layer of application with database and register service.
type App struct {
	GRPCSrv *internal.App
	cfg     *config.Config
}

// New creates upper layer of application
func New(
	log *slog.Logger,
	cfg *config.Config,
	postgresqlClient *pgxpool.Pool,
	sendingClient sending_service.SendingClient,
	jwtClient jwt_service.JWTClient,
) *App {
	log.Info("creating postgresql service")
	storage, err := postgresql.New(postgresqlClient, log)
	if err != nil {
		log.Error("error creating postgresql service: ", err)
		panic(err)
	}

	log.Info("creating auth service")
	authService := auth_service.New(log, cfg, storage)

	log.Info("creating grpc server external_app")
	grpcApp := internal.New(log, cfg, authService, sendingClient, jwtClient)

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
