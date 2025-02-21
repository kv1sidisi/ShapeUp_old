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
	UpdUsr(ctx context.Context, usrmetr *usrdatasvc.UsrMetrics, mask *fieldmaskpb.FieldMask) (updbsusrattr *usrdatasvc.UsrMetrics, err error)
	CreateUsr(ctx context.Context, usrmetr *usrdatasvc.UsrMetrics) (uid []byte, err error)
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

// UpdUsrMetrics is the GRPC server handler method. Updates user's metrics.
//
// Returns:
//
//   - A pointer to UpdMetricsRequest if successful.
//
//   - An error if: Request is invalid. Error while updating user's metrics through service.
func (s *serverAPI) UpdUsrMetrics(ctx context.Context, req *usrdatasvc.UpdUsrMetricsRequest) (*usrdatasvc.UpdUsrMetricsResponse, error) {
	const op = "grpcsrv.UpdBsUsrAttr"
	log := s.log.With(slog.String("op", op))

	if req.GetUser() == nil {
		log.Error("invalid user in grpc request")
		return nil, errdefs.InvalidRequest
	}

	if req.GetUpdMask() == nil || len(req.GetUpdMask().Paths) == 0 {
		log.Error("invalid upd mask in grpc request")
		return nil, errdefs.InvalidRequest
	}

	log.Info("request valid")

	usr, err := s.usrData.UpdUsr(ctx, req.GetUser(), req.GetUpdMask())

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("user attributes updated successfully")

	return &usrdatasvc.UpdUsrMetricsResponse{
		User: usr,
	}, nil
}

// CreateUsrMetrics is the GRPC server handler method. Creates user's metrics.
//
// Returns:
//
//   - A pointer to CreateMetricsRequest if successful.
//
//   - An error if: Request is invalid. Error while creating user's metrics through service.
func (s *serverAPI) CreateUsrMetrics(ctx context.Context, req *usrdatasvc.CreateUsrMetricsRequest) (*usrdatasvc.CreateUsrMetricsResponse, error) {
	const op = "grpcsrv.CreateUsrMetrics"
	log := s.log.With(slog.String("op", op))

	if err := validateCreateUsrMetricsRequest(s.log, req); err != nil {
		return nil, err
	}
	log.Info("request valid")

	uid, err := s.usrData.CreateUsr(ctx, req.GetUser())

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info("user attributes created successfully")

	return &usrdatasvc.CreateUsrMetricsResponse{
		Uid: uid,
	}, nil
}

func validateCreateUsrMetricsRequest(logger *slog.Logger, req *usrdatasvc.CreateUsrMetricsRequest) error {
	const op = "grpcsrv.validateCreateUsrMetricsRequest"
	log := logger.With(slog.String("op", op))

	usr := req.GetUser()

	if usr.Height <= 0 {
		log.Error("invalid user height in grpc request")
		return errdefs.InvalidRequest
	}
	if usr.BirthDate == "" {
		log.Error("invalid user date in grpc request")
		return errdefs.InvalidRequest
	}
	if usr.Gender == "" {
		log.Error("invalid user gender in grpc request")
		return errdefs.InvalidRequest
	}
	if usr.Weight <= 0 {
		log.Error("invalid user weight in grpc request")
		return errdefs.InvalidRequest
	}
	if usr.Name == "" {
		log.Error("invalid user name in grpc request")
		return errdefs.InvalidRequest
	}

	return nil
}
