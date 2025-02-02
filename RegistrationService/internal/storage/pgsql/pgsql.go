package pgsql

import (
	"RegistrationService/internal/storage"
	"RegistrationService/pkg/client/pgsqlcl"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"log/slog"
	"strings"
)

const (
	codeEmailAreadyExists = "23505"
)

type UsrMgr struct {
	client pgsqlcl.Client
	log    *slog.Logger
}

// New creates a new instance of PostreSQL storage
func New(client pgsqlcl.Client, log *slog.Logger) (*UsrMgr, error) {
	return &UsrMgr{
		client: client,
		log:    log,
	}, nil
}

// SaveUser saves user in PostgreSQL database
func (s *UsrMgr) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "pgsql.SaveUser"

	log := s.log.With(
		slog.String("op", op),
	)

	q := `INSERT INTO users (email, password_hash)
			VALUES ($1, $2)
			RETURNING id`

	log.Info(fmt.Sprintf("SQL Query: %s", removeLinesAndTabs(q)))

	if err := s.client.QueryRow(ctx, q, email, passHash).Scan(&uid); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))

			// Error when user with given email already exists
			if pgErr.Code == codeEmailAreadyExists {
				return 0, storage.ErrUserExists
			}

			return 0, newErr
		}
		return 0, err
	}

	return uid, nil
}

func (s *UsrMgr) ConfirmAccount(ctx context.Context, uid int64) error {
	const op = "pgsql.ConfirmAccount"

	log := s.log.With(
		slog.String("op", op),
	)

	query := `UPDATE users SET isconfirmed = TRUE WHERE id = $1`

	log.Info(fmt.Sprintf("SQL Query: %s", query))

	if _, err := s.client.Exec(ctx, query, uid); err != nil {
		return fmt.Errorf("failed to confirm account: %w", err)
	}
	return nil
}

// removeLinesAndTabs removes \n and \t from string.
func removeLinesAndTabs(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
