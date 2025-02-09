package extapp

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kv1sidisi/shapeup/services/authsvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/service/authsvc"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/storage/pgsql"
	"log/slog"
)

// App external layer of GRPC application.
type App struct {
	GRPCSrv *intapp.App
	cfg     *config.Config
}

// New creates services and internal GRPC
//
// Returns App.
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

	authService := authsvc.New(log, cfg, storage, grpccl)
	log.Info("auth service created")

	grpcApp := intapp.New(log, cfg, authService)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
