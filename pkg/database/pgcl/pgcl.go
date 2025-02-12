package pgcl

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	rep "github.com/kv1sidisi/shapeup/pkg/utils/repeatable"
	"log"
	"time"
)

type StorageConfig interface {
	GetDSN() string
}

// Client is interface for connection methods in pgx package.
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// NewClient creates client connection to postgresql throw pgx packet.
func NewClient(ctx context.Context, maxAttempts int, sc StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := sc.GetDSN()
	err = rep.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	return pool, nil
}
