package grpcsrv

import (
	"context"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/api/grpc/pb/usrdatasvc"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
)

// UsrDataSvc service for serverAPI.
type UsrDataSvc interface {
	UpdUsr(ctx context.Context, bsusrattr *usrdatasvc.BsUsrAttr, mask *fieldmaskpb.FieldMask) error
}

// serverAPI handler for the gRPC server.
type serverAPI struct {
	usrdatasvc.UnimplementedUsrDataServer
	usrData UsrDataSvc
	cfg     *config.Config
	log     *slog.Logger
}

// RegisterServer registers services in the GRPC server.
//
// Returns serverAPI as handler for GRPC server.
func RegisterServer(gRPC *grpc.Server,
	usrData UsrDataSvc,
	cfg *config.Config,
	log *slog.Logger,
) {
	usrdatasvc.RegisterUsrDataServer(
		gRPC,
		&serverAPI{
			usrData: usrData,
			cfg:     cfg,
			log:     log,
		})
}
