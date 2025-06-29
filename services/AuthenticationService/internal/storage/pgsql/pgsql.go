package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kv1sidisi/shapeup/pkg/database/pgcl"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/pkg/utils/format"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/domain/models"
	"log/slog"
)

type AuthMgr struct {
	client pgcl.Client
	log    *slog.Logger
}

func New(client pgcl.Client, log *slog.Logger) (*AuthMgr, error) {
	return &AuthMgr{
		client: client,
		log:    log,
	}, nil
}

// FindUserByEmail looks for user model in users database table using email.
//
// Returns:
//   - user model with user's data if successful.
//   - An error if: User not found. Database returns error.
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
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("user not found", err)
			return models.User{}, errdefs.ErrUserNotFound
		}

		log.Error(fmt.Sprintf("query: %s", q), err)
		return models.User{}, errdefs.ErrDatabaseInternal
	}

	return user, nil
}

// AddSession adds new session to sessions database table.
//
// Returns:
//   - userId, access and refresh tokens of successful.
//   - An error if: Session already exists. Database returns error.
func (s *AuthMgr) AddSession(ctx context.Context,
	uid []byte,
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
		log.Error("session already exists: ", uid)
		return errdefs.ErrSessionAlreadyExists
	}

	//TODO: device_info, ip_address
	q := `INSERT INTO user_sessions (uid, refresh_token)
			VALUES ($1, $2)
			RETURNING id`

	log.Info("SQL query: %s", format.RemoveLinesAndTabs(q))

	var sessionId []byte

	if err := s.client.QueryRow(ctx, q, uid, refreshToken).Scan(&sessionId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)
			return errdefs.ErrDatabaseInternal
		}
		log.Error(fmt.Sprintf("query: %s", q), err)
		return errdefs.ErrDatabaseInternal
	}

	log.Info(fmt.Sprintf("new session: %d", sessionId))
	return nil
}

// checkOnlineSession returns true is session with user already exists in session db table.
func checkOnlineSession(ctx context.Context, log *slog.Logger, client pgcl.Client, uid []byte) (bool, error) {
	q := `
        SELECT EXISTS (
            SELECT 1 
            FROM user_sessions 
            WHERE uid = $1
        )`

	log.Info(fmt.Sprintf("query: %s", format.RemoveLinesAndTabs(q)))

	var exists bool
	err := client.QueryRow(ctx, q, uid).Scan(&exists)
	if err != nil {
		log.Error(fmt.Sprintf("query: %s", q), err)
		return false, errdefs.ErrDatabaseInternal
	}
	return exists, nil
}

// IsUserConfirmed checks if user`s account confirmed and returns true or false.
func (s *AuthMgr) IsUserConfirmed(ctx context.Context, uid []byte) (confirmed bool, err error) {
	const op = "postgresql.IsUserConfirmed"
	log := s.log.With(
		slog.String("op", op))

	q := `SELECT is_confirmed FROM users WHERE id = $1`

	log.Info(fmt.Sprintf("query: %s", format.RemoveLinesAndTabs(q)))

	var isConfirmed bool
	if err := s.client.QueryRow(ctx, q, uid).Scan(&isConfirmed); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)
			return false, errdefs.ErrDatabaseInternal
		}
		log.Error(fmt.Sprintf("query: %s", q), err)
		return false, errdefs.ErrDatabaseInternal
	}

	return isConfirmed, nil
}
