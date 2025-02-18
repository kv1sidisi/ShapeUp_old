package extapp

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/app/intapp"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/service/usrdatasvc"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/storage/pgsql"
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
) *App {
	storage, err := pgsql.New(postgresqlClient, log)
	if err != nil {
		panic(err)
	}
	log.Info("postgresql storage manager created")

	usrDataSvc := usrdatasvc.New(log, storage)
	log.Info("user creation service created")

	grpcApp := intapp.New(log, usrDataSvc, cfg)
	log.Info("external GRPC server created")

	return &App{
		GRPCSrv: grpcApp,
		cfg:     cfg,
	}
}
