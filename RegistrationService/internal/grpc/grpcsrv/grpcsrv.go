package grpcsrv

import (
	pbjwt "RegistrationService/api/pb/jwt_service"
	pbsendsvc "RegistrationService/api/pb/sending_service"
	pbusrcreate "RegistrationService/api/pb/user_creation"
	"RegistrationService/internal/config"
	"RegistrationService/internal/storage"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"strings"
)

const (
	confirmAccountLinkBase    = "http://localhost:8082/confirm_account?token="
	confirmationOperationType = "confirmation"
)

// UsrCreationSvc interface represents upper layer of userCreation methods of application.
type UsrCreationSvc interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)

	ConfirmNewUser(
		ctx context.Context,
		jwt string,
		jwtClient pbjwt.JWTClient,
	) (userId int64, err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	pbusrcreate.UnimplementedUserCreationServer
	userCreation  UsrCreationSvc
	cfg           *config.Config
	log           *slog.Logger
	sendingClient pbsendsvc.SendingClient
	jwtClient     pbjwt.JWTClient
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server,
	userCreation UsrCreationSvc,
	cfg *config.Config,
	log *slog.Logger,
	sendingClient pbsendsvc.SendingClient,
	jwtClient pbjwt.JWTClient,
) {
	pbusrcreate.RegisterUserCreationServer(
		gRPC,
		&serverAPI{
			userCreation:  userCreation,
			cfg:           cfg,
			log:           log,
			sendingClient: sendingClient,
			jwtClient:     jwtClient,
		})
}

// Register is the gRPC server handler method, the top layer of the registration process.
func (s *serverAPI) Register(
	ctx context.Context,
	req *pbusrcreate.RegisterRequest,
) (*pbusrcreate.RegisterResponse, error) {
	op := "server.Register"
	log := s.log.With(slog.String("op", op))

	// Validate request in regex
	if err := validateRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Info("register request valid")

	log.Info("registering new user")
	// UsrCreationSvc the new user
	userId, err := s.userCreation.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Info("user registered")

	log.Info("generating confirmation link")
	linkGenResp, err := s.jwtClient.GenerateLink(ctx, &pbjwt.GenerateLinkRequest{
		LinkBase:  confirmAccountLinkBase,
		Uid:       userId,
		Operation: confirmationOperationType,
	})
	if err != nil {
		log.Error("confirmation link generation failed")
		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Info("confirmation link generated successfully: ", linkGenResp.GetLink())

	log.Info("sending user confirmation link")
	sendEmailResp, err := s.sendingClient.SendEmail(ctx, &pbsendsvc.EmailRequest{
		Message: linkGenResp.GetLink(),
		Email:   req.GetEmail(),
	})
	if err != nil {
		log.Error("failed to send confirmation link")
		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Info("user confirmation link sent successfully to:" + sendEmailResp.GetEmail())

	// Return the response with the user ID
	return &pbusrcreate.RegisterResponse{
		UserId: userId,
	}, nil
}

// validateRegisterRequest performs validation on the registration request.
func validateRegisterRequest(req *pbusrcreate.RegisterRequest) error {
	// Validate email
	if !govalidator.IsEmail(req.GetEmail()) {
		return errors.New("incorrect email address")
	}

	// Validate password length
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Validate spaces in password
	if strings.Contains(req.GetPassword(), " ") {
		return errors.New("password must not contain spaces")
	}

	return nil
}

// Confirm is the gRPC server handler method, the top layer of the registration process.
func (s *serverAPI) Confirm(ctx context.Context,
	req *pbusrcreate.ConfirmRequest,
) (*pbusrcreate.ConfirmResponse, error) {

	s.log.Info("confirming new user with token")
	userId, err := s.userCreation.ConfirmNewUser(ctx, req.Jwt, s.jwtClient)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	s.log.Info("user confirmed")

	return &pbusrcreate.ConfirmResponse{
		UserId: userId,
	}, nil
}
