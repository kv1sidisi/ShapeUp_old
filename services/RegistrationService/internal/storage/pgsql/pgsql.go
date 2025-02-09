package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/kv1sidisi/shapeup/libs/common/errdefs"
	"github.com/kv1sidisi/shapeup/services/regsvc/pkg/client/pgsqlcl"
	"log/slog"
	"strings"
)

const (
	codeEmailAlreadyExists = "23505"
)

type UsrMgr struct {
	client pgsqlcl.Client
	log    *slog.Logger
}

func New(client pgsqlcl.Client, log *slog.Logger) (*UsrMgr, error) {
	return &UsrMgr{
		client: client,
		log:    log,
	}, nil
}

// SaveUser saves user in PostgreSQL database.
//
// Returns:
//   - user ID if successful.
//   - An error if: Email already exists. Database returns error.
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
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)

			if pgErr.Code == codeEmailAlreadyExists {
				log.Error("email already exists")
				return 0, errdefs.ErrEmailAlreadyExists
			}

			return 0, errdefs.ErrDatabaseInternal
		}
		return 0, errdefs.ErrDatabaseInternal
	}

	return uid, nil
}

// ConfirmAccount confirms account in PostgreSQL database.
//
// Returns:
//   - An error if: Database returns error.
func (s *UsrMgr) ConfirmAccount(ctx context.Context, uid int64) error {
	const op = "pgsql.ConfirmAccount"
	log := s.log.With(
		slog.String("op", op),
	)

	query := `UPDATE users SET isconfirmed = TRUE WHERE id = $1`

	log.Info(fmt.Sprintf("SQL Query: %s", query))

	if _, err := s.client.Exec(ctx, query, uid); err != nil {
		log.Error("failed to confirm account: ", err)
		return errdefs.ErrDatabaseInternal
	}
	return nil
}

// DeleteUser deletes account in PostgreSQL database.
//
// Returns:
//   - An error if: Database returns error.
func (s *UsrMgr) DeleteUser(ctx context.Context, uid int64) error {
	const op = "pgsql.DeleteUser"
	log := s.log.With(
		slog.String("op", op),
	)
	q := `DELETE FROM users WHERE id = $1`

	log.Info(fmt.Sprintf("SQL query: %s", removeLinesAndTabs(q)))

	if _, err := s.client.Exec(ctx, q, uid); err != nil {
		log.Error("failed to delete user: ", err)
		return errdefs.ErrDatabaseInternal
	}
	return nil

}

// removeLinesAndTabs removes \n and \t from string.
func removeLinesAndTabs(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
