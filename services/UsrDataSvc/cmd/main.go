package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	loadconfig "github.com/kv1sidisi/shapeup/pkg/config"
	"github.com/kv1sidisi/shapeup/pkg/database/pgcl"
	"github.com/kv1sidisi/shapeup/pkg/logger"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/config"
	"log/slog"
)

func main() {
	cfg := &config.Config{}
	loadconfig.MustLoad(cfg)

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting up", slog.String("env", cfg.Env))

	postgresqlClient := mustConnectToDatabase(cfg)
	log.Info("connected to database")

	application

}

// mustConnectToDatabase panics if setupDatabaseConnection fails
func mustConnectToDatabase(cfg *config.Config) *pgxpool.Pool {
	postgresqlClient, err := setupDatabaseConnection(cfg)
	if err != nil {
		panic(err)
	}
	return postgresqlClient
}

// setupDatabaseConnection connect to database
// returns pgx client.
func setupDatabaseConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	postgresqlClient, err := pgcl.NewClient(context.TODO(), 3, &cfg.Storage)
	if err != nil {
		return nil, err
	}
	return postgresqlClient, nil
}
