package usrdatasvc

import (
	"context"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/api/grpc/pb/usrdatasvc"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/domain/models"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
)

type UsrDataSvc struct {
	log        *slog.Logger
	usrDataMgr UsrDataMgr
}

type UsrDataMgr interface {
	UpdBsUsrAttr(ctx context.Context, bsusrattr *models.BsUsrAttr) error
	GetById(ctx context.Context, id int64) (*models.BsUsrAttr, error)
}

func New(log *slog.Logger, usrDataMgr UsrDataMgr) *UsrDataSvc {
	return &UsrDataSvc{log, usrDataMgr}
}

func (u *UsrDataSvc) UpdUsr(ctx context.Context, bsusrattr *usrdatasvc.BsUsrAttr, mask *fieldmaskpb.FieldMask) error {
	panic("implement me")
}
