package postgresql

import (
	"AuthenticationService/pkg/client/postgresql"
	"context"
	"log/slog"
)

type Storage struct {
	client postgresql.Client
	log    *slog.Logger
}

func New(client postgresql.Client, log *slog.Logger) (*Storage, error) {
	return &Storage{
		client: client,
		log:    log,
	}, nil
}

func (s *Storage) FindUser(ctx context.Context,
	username string,
	password string,
) (uid int64, err error) {
	return 1, nil
}
