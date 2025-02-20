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
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/domain/models"
	"log/slog"
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

func (s *UsrDataMgr) UpdBsUsrAttr(ctx context.Context, bsusrattr *models.BsUsrAttr) error {
	panic("implement me")
}

func (s *UsrDataMgr) GetById(ctx context.Context) (*models.BsUsrAttr, error) {
	const op = "postgresql.GetById"
	log := s.log.With(
		slog.String("op", op))

	q := `SELECT name, height, weight, gender, birth_date FROM bsusrattr WHERE uid = $1`

	log.Info(fmt.Sprintf("query: %s", format.RemoveLinesAndTabs(q)))

	usrdata := &models.BsUsrAttr{}
	uid := ctx.Value("uid")

	if err := s.client.QueryRow(ctx, q, uid).Scan(
		&usrdata.Name, &usrdata.Height, &usrdata.Weight, &usrdata.Gender, &usrdata.BirthDate,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()), err)
			return nil, errdefs.ErrDatabaseInternal
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Profile not found, uid: ", uid)
			return nil, nil
		}
		log.Error(fmt.Sprintf("query: %s", q), err)
		return nil, errdefs.ErrDatabaseInternal
	}

	return usrdata, nil

}
