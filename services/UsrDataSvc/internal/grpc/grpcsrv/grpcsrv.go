package grpcsrv

import (
	"context"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/pkg/proto/usrdatasvc/pb"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
)

// UsrDataSvc service for serverAPI.
type UsrDataSvc interface {
	UpdUsr(ctx context.Context, bsusrattr *usrdatasvc.BsUsrAttr, mask *fieldmaskpb.FieldMask) (updbsusrattr *usrdatasvc.BsUsrAttr, err error)
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

// UpdBsUsrAttr is the GRPC server handler method. Updates base user's attributes.
//
// Returns:
//
//   - A pointer to UpdBsUsrAttrRequest if successful.
//
//   - An error if: Request is invalid. Error while updating base user's attributes through service.
func (s *serverAPI) UpdBsUsrAttr(ctx context.Context, req *usrdatasvc.UpdBsUsrAttrRequest) (*usrdatasvc.UpdBsUsrAttrResponse, error) {
	const op = "grpcsrv.UpdBsUsrAttr"
	log := s.log.With(slog.String("op", op))

	if req.GetUser() == nil {
		log.Error("invalid user in grpc request")
		return nil, errdefs.InvalidRequest
	}

	if ctx.Value("uid") == nil {
		log.Error("invalid uid in grpc request(context)")
		return nil, errdefs.InvalidCredentials
	}

	log.Info("request valid")

	updbsusrattr, err := s.usrData.UpdUsr(ctx, req.GetUser(), req.GetUpdMask())

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("user attributes updated successfully")

	return &usrdatasvc.UpdBsUsrAttrResponse{
		User: updbsusrattr,
	}, nil
}
