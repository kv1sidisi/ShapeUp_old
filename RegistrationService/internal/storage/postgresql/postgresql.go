package postgresql

import (
	"RegistrationService/internal/storage"
	"RegistrationService/pkg/client/postgresql"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"log/slog"
)

type Storage struct {
	client postgresql.Client
	log    *slog.Logger
}

// New creates a new instance of PostreSQL storage
func New(client postgresql.Client, log *slog.Logger) (*Storage, error) {
	return &Storage{
		client: client,
		log:    log,
	}, nil
}

// SaveUser saves user in PostgreSQL database
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "postgresql.postgresql.SaveUser"

	log := s.log.With(
		slog.String("op", op),
	)

	q := `insert into users (email, password_hash)
			values ($1, $2)
			returning id`

	log.Info(fmt.Sprintf("SQL Query: %s", q))

	if err := s.client.QueryRow(ctx, q, email, passHash).Scan(&uid); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Error(newErr.Error())

			// Error when user with given email already exists
			if pgErr.Code == "23505" {
				return 0, storage.ErrUserExists
			}

			return 0, nil
		}
		return 0, err
	}

	return uid, nil
}

func (s *Storage) ConfirmAccount(ctx context.Context, uid int64) error {
	const op = "storage.postgresql.ConfirmAccount"

	log := s.log.With(
		slog.String("op", op),
	)

	query := `UPDATE users SET isconfirmed = TRUE WHERE id = $1`

	log.Info(fmt.Sprintf("SQL Query: %s", query))

	if _, err := s.client.Exec(ctx, query, uid); err != nil {
		log.Error("failed to confirm account", err)
		return fmt.Errorf("failed to confirm account: %w", err)
	}
	return nil
}
