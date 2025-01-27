package grpc_server

import (
	sendv1 "SendingService/api/pb"
	"SendingService/internal/config"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"log/slog"
)

type SendingService interface {
	SendEmail(
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
	op := "server.SendEmail"

	log := s.log.With(slog.String("op", op))

	log.Info("validating email")
	if !govalidator.IsEmail(req.GetEmail()) {
		return nil, errors.New("incorrect email address: " + req.GetEmail())
	}
	log.Info("email valid")

	log.Info("sending email")
	//TODO: send email service initialization
	log.Info("email sent")

	return &sendv1.EmailResponse{
		Email: req.GetEmail(),
	}, nil
}
