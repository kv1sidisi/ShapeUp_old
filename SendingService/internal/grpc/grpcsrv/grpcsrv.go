package grpcsrv

import (
	"SendingService/api/pb/sending_service"
	"SendingService/internal/config"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// SendingService represent service for sending, bottom layer.
type SendingService interface {
	SMTPSendNewEmail(
		ctx context.Context,
		email string,
		message string,
	) error
	GoGetSendNewEmail(
		ctx context.Context,
		email string,
		message string,
	) error
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	sending_service.UnimplementedSendingServer
	sendingService SendingService
	cfg            *config.Config
	log            *slog.Logger
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, sendingService SendingService, cfg *config.Config, log *slog.Logger) {
	sending_service.RegisterSendingServer(gRPC,
		&serverAPI{
			sendingService: sendingService,
			cfg:            cfg,
			log:            log,
		})
}

// SendEmail is the gRPC server handler method, the top layer of the sending process.
func (s *serverAPI) SendEmail(
	ctx context.Context,
	req *sending_service.EmailRequest,
) (*sending_service.EmailResponse, error) {
	const op = "server.SendEmail"

	log := s.log.With(slog.String("op", op))

	log.Info("validating email")
	if !govalidator.IsEmail(req.GetEmail()) {
		return nil, errors.New("incorrect email address: " + req.GetEmail())
	}
	log.Info("email valid")

	log.Info("sending email")
	if err := s.sendingService.GoGetSendNewEmail(ctx, req.GetEmail(), req.GetMessage()); err != nil {
		log.Error("sending email error:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Info("email sent")

	return &sending_service.EmailResponse{
		Email: req.GetEmail(),
	}, nil
}
