package grpcsrv

import (
	usrcreatesvc2 "RegistrationService/api/grpc/pb/usrcreatesvc"
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

// UsrCreateSvc interface of user creation service.
type UsrCreateSvc interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)

	ConfirmNewUser(
		ctx context.Context,
		jwt string,
	) (userId int64, err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	usrcreatesvc2.UnimplementedUserCreationServer
	userCreation UsrCreateSvc
	cfg          *config.Config
	log          *slog.Logger
}

// RegisterServer registers the request handler in the gRPC server.
func RegisterServer(gRPC *grpc.Server,
	userCreation UsrCreateSvc,
	cfg *config.Config,
	log *slog.Logger,
) {
	usrcreatesvc2.RegisterUserCreationServer(
		gRPC,
		&serverAPI{
			userCreation: userCreation,
			cfg:          cfg,
			log:          log,
		})
}

// Register is the gRPC server handler method, the top layer of the registration process.
func (s *serverAPI) Register(
	ctx context.Context,
	req *usrcreatesvc2.RegisterRequest,
) (*usrcreatesvc2.RegisterResponse, error) {
	const op = "grpcsrv.Register"
	log := s.log.With(slog.String("op", op))

	// Validate request in regex
	if err := validateRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log.Info("register request valid")

	// UsrCreateSvc the new user
	userId, err := s.userCreation.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}
	log.Info("user registered successfully")

	// Return the response with the user ID
	return &usrcreatesvc2.RegisterResponse{
		UserId: userId,
	}, nil
}

// validateRegisterRequest performs validation on the registration request.
func validateRegisterRequest(req *usrcreatesvc2.RegisterRequest) error {
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
	req *usrcreatesvc2.ConfirmRequest,
) (*usrcreatesvc2.ConfirmResponse, error) {

	userId, err := s.userCreation.ConfirmNewUser(ctx, req.Jwt)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	s.log.Info("user confirmed")

	return &usrcreatesvc2.ConfirmResponse{
		UserId: userId,
	}, nil
}
