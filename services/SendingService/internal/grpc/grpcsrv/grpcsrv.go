package grpcsrv

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/kv1sidisi/shapeup/libs/common/errdefs"
	"github.com/kv1sidisi/shapeup/services/sendsvc/api/grpc/pb/sendsvc"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/config"
	"google.golang.org/grpc"
	"log/slog"
)

// SendSvc represent service for sending, bottom layer.
type SendSvc interface {
	GoGetSendNewEmail(
		ctx context.Context,
		email string,
		message string,
	) error
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	sendsvc.UnimplementedSendingServer
	sendingService SendSvc
	cfg            *config.Config
	log            *slog.Logger
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, sendingService SendSvc, cfg *config.Config, log *slog.Logger) {
	sendsvc.RegisterSendingServer(gRPC,
		&serverAPI{
			sendingService: sendingService,
			cfg:            cfg,
			log:            log,
		})
}

// SendEmail is the gRPC server handler method, the top layer of the sending process.
func (s *serverAPI) SendEmail(
	ctx context.Context,
	req *sendsvc.EmailRequest,
) (*sendsvc.EmailResponse, error) {
	const op = "server.SendEmail"
	log := s.log.With(slog.String("op", op))

	if !govalidator.IsEmail(req.GetEmail()) {
		log.Error("invalid email")
		return nil, errdefs.InvalidEmail
	}
	log.Info("email valid")

	if err := s.sendingService.GoGetSendNewEmail(ctx, req.GetEmail(), req.GetMessage()); err != nil {
		return nil, err
	}

	return &sendsvc.EmailResponse{
		Email: req.GetEmail(),
	}, nil
}
