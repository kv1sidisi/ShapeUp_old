package extapp

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kv1sidisi/shapeup/services/regsvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/service/usrcreatesvc"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/storage/pgsql"
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

	usrCreateSvc := usrcreatesvc.New(log, storage, grpccl)
	log.Info("user creation service created")

	grpcApp := intapp.New(log, usrCreateSvc, cfg)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
