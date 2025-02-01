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
)

const (
	ErrUserNotFound          = "user not found"
	ErrSessionAlreadyExists  = "session already exists"
	codeSessionAlreadyExists = "23505"
)

type Storage struct {
	client pgsqlcl.Client
	log    *slog.Logger
}

func New(client pgsqlcl.Client, log *slog.Logger) (*Storage, error) {
	return &Storage{
		client: client,
		log:    log,
	}, nil
}

func (s *Storage) FindUserByEmail(ctx context.Context,
	email string,
) (user models.User, err error) {
	const op = "postgresql.FindUserByEmail"

	log := s.log.With(
		slog.String("op", op))

	q := `SELECT id, email, password_hash FROM users WHERE email = $1`

	log.Info(fmt.Sprintf("Query: %s", q))

	if err := s.client.QueryRow(ctx, q, email).Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) AddSession(ctx context.Context,
	uid int64,
	accessToken string,
	refreshToken string,
) (err error) {
	const op = "postgresql.SaveSession"

	log := s.log.With(
		"op", op,
	)

	q := `INSERT INTO sessions (user_id, access_token, refresh_token)
			VALUES ($1, $2, $3)
			RETURNING id`

	log.Info("SQL Query: %s", q)

	var sessionId int64

	if err := s.client.QueryRow(ctx, q, uid, accessToken, refreshToken).Scan(&sessionId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Info(newErr.Error())

			if pgErr.Code == codeSessionAlreadyExists {
				return fmt.Errorf(ErrSessionAlreadyExists)
			}

			return nil
		}
		return err
	}

	log.Info(fmt.Sprintf("new session: %d", sessionId))
	return nil
}
