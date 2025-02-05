package grpcsrv

import (
	"SendingService/api/grpc/pb/sendsvc"
	"SendingService/internal/config"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, errors.New("incorrect email address: " + req.GetEmail())
	}
	log.Info("email valid")

	if err := s.sendingService.GoGetSendNewEmail(ctx, req.GetEmail(), req.GetMessage()); err != nil {
		log.Error("sending email error:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sendsvc.EmailResponse{
		Email: req.GetEmail(),
	}, nil
}
