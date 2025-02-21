package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/kv1sidisi/shapeup/pkg/database/pgcl"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/pkg/utils/format"
	"log/slog"
)

const (
	codeEmailAlreadyExists = "23505"
)

type UsrMgr struct {
	client pgcl.Client
	log    *slog.Logger
}

func New(client pgcl.Client, log *slog.Logger) (*UsrMgr, error) {
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
func (s *UsrMgr) SaveUser(ctx context.Context, email string, passHash []byte) (uid []byte, err error) {
	const op = "pgsql.SaveUser"

	log := s.log.With(
		slog.String("op", op),
	)

	q := `INSERT INTO users (email, password_hash)
			VALUES ($1, $2)
			RETURNING id`

	log.Info(fmt.Sprintf("SQL Query: %s", format.RemoveLinesAndTabs(q)))

	if err := s.client.QueryRow(ctx, q, email, passHash).Scan(&uid); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)

			if pgErr.Code == codeEmailAlreadyExists {
				log.Error("email already exists")
				return uid, errdefs.ErrEmailAlreadyExists
			}

			return uid, errdefs.ErrDatabaseInternal
		}
		log.Error(err.Error())
		return uid, errdefs.ErrDatabaseInternal
	}

	return uid, nil
}

// ConfirmAccount confirms account in PostgreSQL database.
//
// Returns:
//   - An error if: Database returns error.
func (s *UsrMgr) ConfirmAccount(ctx context.Context, uid []byte) error {
	const op = "pgsql.ConfirmAccount"
	log := s.log.With(
		slog.String("op", op),
	)

	query := `UPDATE users SET is_confirmed = TRUE WHERE id = $1`

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
func (s *UsrMgr) DeleteUser(ctx context.Context, uid []byte) error {
	const op = "pgsql.DeleteUser"
	log := s.log.With(
		slog.String("op", op),
	)
	q := `DELETE FROM users WHERE id = $1`

	log.Info(fmt.Sprintf("SQL query: %s", format.RemoveLinesAndTabs(q)))

	if _, err := s.client.Exec(ctx, q, uid); err != nil {
		log.Error("failed to delete user: ", err)
		return errdefs.ErrDatabaseInternal
	}
	return nil

}
