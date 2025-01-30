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

func (s *Storage) FindUserByName(ctx context.Context,
	username string,
	password string,
) (uid int64, err error) {
	return 1, nil
}

func (s *Storage) FindUserByEmail(ctx context.Context,
	username string,
	password string,
) (uid int64, err error) {
	return 1, nil
}

func (s *Storage) AddSession(ctx context.Context,
	uid int64,
	accessToken string,
	refreshToken string,
) (err error) {
	return nil
}
