package grpc_server

import (
	sendv1 "SendingService/api/pb"
	"SendingService/internal/config"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type SendingService interface {
	SendNewEmail(
		ctx context.Context,
		email string,
		message string,
	) error
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	sendv1.UnimplementedSendingServer
	sendingService SendingService
	cfg            *config.Config
	log            *slog.Logger
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server, sendingService SendingService, cfg *config.Config, log *slog.Logger) {
	sendv1.RegisterSendingServer(gRPC,
		&serverAPI{
			sendingService: sendingService,
			cfg:            cfg,
			log:            log,
		})
}

func (s *serverAPI) SendEmail(
	ctx context.Context,
	req *sendv1.EmailRequest,
) (*sendv1.EmailResponse, error) {
	const op = "server.SendEmail"

	log := s.log.With(slog.String("op", op))

	log.Info("validating email")
	if !govalidator.IsEmail(req.GetEmail()) {
		return nil, errors.New("incorrect email address: " + req.GetEmail())
	}
	log.Info("email valid")

	log.Info("sending email")
	if err := s.sendingService.SendNewEmail(ctx, req.GetEmail(), req.GetMessage()); err != nil {
		log.Error("sending email error:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Info("email sent")

	return &sendv1.EmailResponse{
		Email: req.GetEmail(),
	}, nil
}
