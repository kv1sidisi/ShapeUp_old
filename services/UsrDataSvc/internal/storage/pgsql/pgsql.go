package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kv1sidisi/shapeup/pkg/database/pgcl"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/pkg/utils/format"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/domain/models"
	"log/slog"
	"time"
)

const (
	codeUserMetricsAlreadyExists = "23505"
)

type UsrDataMgr struct {
	client pgcl.Client
	log    *slog.Logger
}

func New(client pgcl.Client, log *slog.Logger) (*UsrDataMgr, error) {
	return &UsrDataMgr{
		client: client,
		log:    log,
	}, nil
}

// UpdUsrMetrics saves user in PostgreSQL database.
//
// Returns:
//   - nil if successful.
//   - An error if: User metrics not found. Database returns error.
func (s *UsrDataMgr) UpdUsrMetrics(ctx context.Context, usrmetr *models.UsrMetrics) error {
	const op = "postgresql.UpdBsUsrAttr"
	log := s.log.With(slog.String("op", op))

	q := `UPDATE user_metrics
          SET name = $2, height = $3, weight = $4, birth_date = $5, gender = $6
          WHERE uid = $1`

	log.Info(fmt.Sprintf("SQL Query: %s", format.RemoveLinesAndTabs(q)))

	uid := ctx.Value("uid").([]byte)
	if uid == nil {
		log.Error("uid is nil in context")
		return errdefs.ErrInternal
	}

	if _, err := s.client.Exec(ctx, q, uid, usrmetr.Name, usrmetr.Height, usrmetr.Weight, usrmetr.BirthDate, usrmetr.Gender); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("user metrics not found, uid: %x", uid)
			return errdefs.ErrUserNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)

			return errdefs.ErrDatabaseInternal
		}
		return errdefs.ErrDatabaseInternal
	}
	return nil
}

// GetById saves user in PostgreSQL database.
//
// Returns:
//   - User metrics data if successful.
//   - An error if: User does not exist. Database returns error.
func (s *UsrDataMgr) GetById(ctx context.Context) (*models.UsrMetrics, error) {
	const op = "postgresql.GetById"
	log := s.log.With(slog.String("op", op))

	q := `SELECT name, height, weight, birth_date, gender FROM user_metrics WHERE uid = $1`

	log.Info(fmt.Sprintf("query: %s", format.RemoveLinesAndTabs(q)))
	usrdata := &models.UsrMetrics{}

	rawUid, ok := ctx.Value("uid").([]byte)
	if !ok {
		log.Error("uid is not []byte")
		return nil, errdefs.ErrInternal
	}
	if rawUid == nil {
		log.Error("uid is nil in context")
		return nil, errdefs.ErrInternal
	}
	uid, err := uuid.FromBytes(rawUid)
	if err != nil {
		log.Error("failed to convert uid: %v", err)
		return nil, errdefs.ErrInternal
	}

	log.Info(fmt.Sprintf("uid: %s", uid))

	var birthDate time.Time
	if err := s.client.QueryRow(ctx, q, uid).Scan(&usrdata.Name, &usrdata.Height, &usrdata.Weight, &birthDate, &usrdata.Gender); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Profile not found, uid: %x", uid)
			return nil, errdefs.ErrUserNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil, errdefs.ErrDatabaseInternal
		}

		log.Error("Query failed: %s, Error: %v", q, err)
		return nil, errdefs.ErrDatabaseInternal
	}

	usrdata.BirthDate = birthDate.Format("2006-01-02")
	return usrdata, nil
}

// CreateUsrMetrics saves user in PostgreSQL database.
//
// Returns:
//   - uid if successful.
//   - An error if: User metrics already exists. Database returns error.
func (s *UsrDataMgr) CreateUsrMetrics(ctx context.Context, usrmetr *models.UsrMetrics) (uid []byte, err error) {
	const op = "postgresql.CreateUsrMetrics"
	log := s.log.With(slog.String("op", op))

	uid = ctx.Value("uid").([]byte)
	if uid == nil {
		log.Error("uid is nil in context")
		return nil, errdefs.ErrInternal
	}

	q := `INSERT INTO user_metrics (uid, name, height, weight, birth_date, gender)
			VALUES ($1, $2, $3, $4, $5, $6)`

	log.Info(fmt.Sprintf("SQL Query: %s", format.RemoveLinesAndTabs(q)))

	if _, err := s.client.Exec(ctx, q, uid, usrmetr.Name, usrmetr.Height, usrmetr.Weight, usrmetr.BirthDate, usrmetr.Gender); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)

			if pgErr.Code == codeUserMetricsAlreadyExists {
				log.Error("user metrics already exists")
				return uid, errdefs.ErrUserMetricsAlreadyExists
			}

			return uid, errdefs.ErrDatabaseInternal
		}
		log.Error(err.Error())
		return uid, errdefs.ErrDatabaseInternal
	}

	return uid, nil
}
