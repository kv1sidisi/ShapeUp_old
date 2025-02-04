package postgresql

import (
	"AuthenticationService/internal/domain/models"
	"AuthenticationService/pkg/client/pgsqlcl"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"log/slog"
	"strings"
)

const (
	ErrUserNotFound         = "user not found"
	ErrSessionAlreadyExists = "session already exists"
)

type AuthMgr struct {
	client pgsqlcl.Client
	log    *slog.Logger
}

func New(client pgsqlcl.Client, log *slog.Logger) (*AuthMgr, error) {
	return &AuthMgr{
		client: client,
		log:    log,
	}, nil
}

func (s *AuthMgr) FindUserByEmail(ctx context.Context,
	email string,
) (user models.User, err error) {
	const op = "postgresql.FindUserByEmail"

	log := s.log.With(
		slog.String("op", op))

	q := `SELECT id, email, password_hash FROM users WHERE email = $1`

	log.Info(fmt.Sprintf("query: %s", q))

	if err := s.client.QueryRow(ctx, q, email).Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %s", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *AuthMgr) AddSession(ctx context.Context,
	uid int64,
	accessToken string,
	refreshToken string,
) (err error) {
	const op = "postgresql.SaveSession"

	log := s.log.With(
		"op", op,
	)

	exists, err := checkOnlineSession(ctx, log, s.client, uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return fmt.Errorf("%s: %w", op, ErrSessionAlreadyExists)
	}

	q := `INSERT INTO sessions (user_id, access_token, refresh_token)
			VALUES ($1, $2, $3)
			RETURNING id`

	log.Info("SQL query: %s", removeLinesAndTabs(q))

	var sessionId int64

	if err := s.client.QueryRow(ctx, q, uid, accessToken, refreshToken).Scan(&sessionId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		}
		return err
	}

	log.Info(fmt.Sprintf("new session: %d", sessionId))
	return nil
}

func checkOnlineSession(ctx context.Context, log *slog.Logger, client pgsqlcl.Client, uid int64) (bool, error) {
	q := `
        SELECT EXISTS (
            SELECT 1 
            FROM sessions 
            WHERE user_id = $1
        )`

	log.Info(fmt.Sprintf("query: %s", removeLinesAndTabs(q)))

	var exists bool
	err := client.QueryRow(ctx, q, uid).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query error: %w", err)
	}
	return exists, nil
}

// IsUserConfirmed checks if user`s account confirmed and returns true or false.
func (s *AuthMgr) IsUserConfirmed(ctx context.Context, uid int64) (confirmed bool, err error) {
	const op = "postgresql.IsUserConfirmed"
	log := s.log.With(
		slog.String("op", op))

	q := `SELECT isconfirmed FROM users WHERE id = $1`

	log.Info(fmt.Sprintf("query: %s", removeLinesAndTabs(q)))

	var isConfirmed bool
	if err := s.client.QueryRow(ctx, q, uid).Scan(&isConfirmed); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return false, fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		}
		return false, err
	}

	return isConfirmed, nil
}

// removeLinesAndTabs removes \n and \t from string.
func removeLinesAndTabs(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
