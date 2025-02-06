package pgsql

import (
	"AuthenticationService/internal/domain/models"
	"AuthenticationService/pkg/client/pgsqlcl"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// FindUserByEmail looks for user model in users db table using email.
func (s *AuthMgr) FindUserByEmail(ctx context.Context,
	email string,
) (user models.User, err error) {
	const op = "postgresql.FindUserByEmail"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email))

	q := `SELECT id, email, password_hash FROM users WHERE email = $1`

	log.Info(fmt.Sprintf("query: %s", q))

	if err := s.client.QueryRow(ctx, q, email).Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, status.Error(codes.NotFound, ErrUserNotFound)
		}

		return models.User{}, status.Error(codes.Internal, err.Error())
	}

	return user, nil
}

// AddSession adds new session to sessions db table.
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
		return err
	}
	if exists {
		return status.Error(codes.AlreadyExists, ErrSessionAlreadyExists)
	}

	q := `INSERT INTO sessions (user_id, access_token, refresh_token)
			VALUES ($1, $2, $3)
			RETURNING id`

	log.Info("SQL query: %s", removeLinesAndTabs(q))

	var sessionId int64

	if err := s.client.QueryRow(ctx, q, uid, accessToken, refreshToken).Scan(&sessionId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return status.Error(codes.Internal, fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Info(fmt.Sprintf("new session: %d", sessionId))
	return nil
}

// checkOnlineSession returns true is session with user already exists in session db table.
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
		return false, status.Error(codes.Internal, err.Error())
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
			return false, status.Error(codes.Internal, fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		}
		return false, status.Error(codes.Internal, err.Error())
	}

	return isConfirmed, nil
}

// removeLinesAndTabs removes \n and \t from string.
func removeLinesAndTabs(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
