package pgsql

import (
	"context"
	"github.com/kv1sidisi/shapeup/pkg/database/pgcl"
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
