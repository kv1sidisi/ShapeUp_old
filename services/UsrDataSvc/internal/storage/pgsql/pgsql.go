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

func (s *UsrDataMgr) UpdUsrMetrics(ctx context.Context, usrmetr *models.UsrMetrics) error {
	const op = "postgresql.UpdBsUsrAttr"
	log := s.log.With(slog.String("op", op))

	q := `INSERT INTO user_metrics (uid, name, height, weight, birth_date, gender)
		VALUES ($1, $2, $3, $4, $5, $6)`

	log.Info(fmt.Sprintf("SQL Query: %s", format.RemoveLinesAndTabs(q)))

	uid := ctx.Value("uid")
	if err := s.client.QueryRow(ctx, q, uid, usrmetr.Name, usrmetr.Height, usrmetr.Weight, usrmetr.BirthDate, usrmetr.Gender).Scan(&uid); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)

			if pgErr.Code == codeUserMetricsAlreadyExists {
				log.Error("email already exists")
				return errdefs.ErrUserMetricsAlreadyExists
			}

			return errdefs.ErrDatabaseInternal
		}
		return errdefs.ErrDatabaseInternal
	}

	return nil
}

func (s *UsrDataMgr) GetById(ctx context.Context) (*models.UsrMetrics, error) {
	const op = "postgresql.GetById"
	log := s.log.With(slog.String("op", op))

	q := `SELECT name, height, weight, birth_date, gender FROM user_metrics WHERE uid = $1`

	log.Info(fmt.Sprintf("query: %s", format.RemoveLinesAndTabs(q)))
	usrdata := &models.UsrMetrics{}

	rawUID := ctx.Value("uid")
	if rawUID == nil {
		log.Error("uid is nil in context")
		return nil, errdefs.ErrInternal
	}

	uid, err := uuid.FromBytes(rawUID.([]byte))
	if err != nil {
		log.Error("Failed to convert UID: %v", err)
		return nil, errdefs.ErrInternal
	}

	log.Info("Fetching profile with uid: %x", uid)

	err = s.client.QueryRow(ctx, q, uid.String()).Scan(
		&usrdata.Name, &usrdata.Height, &usrdata.Weight, &usrdata.Gender, &usrdata.BirthDate,
	)
	if err != nil {
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

	return usrdata, nil
}
