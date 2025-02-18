package grpcsrv

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	sendsvc "github.com/kv1sidisi/shapeup/pkg/proto/sendsvc/pb"
	"github.com/kv1sidisi/shapeup/services/sendsvc/internal/config"
	"google.golang.org/grpc"
	"log/slog"
)

// SendSvc service for serverAPI.
type SendSvc interface {
	GoGetSendNewEmail(
		ctx context.Context,
		email string,
		message string,
	) error
}

// serverAPI handler for the GRPC server.
type serverAPI struct {
	sendsvc.UnimplementedSendingServer
	sendingService SendSvc
	cfg            *config.Config
	log            *slog.Logger
}

// RegisterServer registers services in the GRPC server.
//
// Returns serverAPI as handler for GRPC server.
func RegisterServer(gRPC *grpc.Server, sendingService SendSvc, cfg *config.Config, log *slog.Logger) {
	sendsvc.RegisterSendingServer(gRPC,
		&serverAPI{
			sendingService: sendingService,
			cfg:            cfg,
			log:            log,
		})
}

// SendEmail is the GRPC server handler method. Sends email.
//
// Returns:
//
//   - A pointer to EmailResponse if successful.
//
//   - An error if: Email is invalid. Error while sending the email.
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
