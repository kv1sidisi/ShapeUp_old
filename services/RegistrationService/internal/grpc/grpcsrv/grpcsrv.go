package grpcsrv

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/kv1sidisi/shapeup/libs/common/errdefs"
	usrcreatesvc2 "github.com/kv1sidisi/shapeup/services/regsvc/api/grpc/pb/usrcreatesvc"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/config"
	"google.golang.org/grpc"
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
	if err := validateRegisterRequest(log, req); err != nil {
		return nil, err
	}
	log.Info("register request valid")

	// UsrCreateSvc the new user
	userId, err := s.userCreation.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	log.Info("user registered successfully")

	// Return the response with the user ID
	return &usrcreatesvc2.RegisterResponse{
		UserId: userId,
	}, nil
}

// validateRegisterRequest performs validation on the registration request.
func validateRegisterRequest(log *slog.Logger, req *usrcreatesvc2.RegisterRequest) error {
	// Validate email
	if !govalidator.IsEmail(req.GetEmail()) {
		log.Error("invalid email in request")
		return errdefs.InvalidCredentials
	}

	// Validate password length
	if len(req.Password) < 8 {
		log.Error("invalid password in request")
		return errdefs.InvalidCredentials
	}

	// Validate spaces in password
	if strings.Contains(req.GetPassword(), " ") {
		log.Error("invalid password in request")
		return errdefs.InvalidCredentials
	}

	return nil
}

// Confirm is the gRPC server handler method, the top layer of the registration process.
func (s *serverAPI) Confirm(ctx context.Context,
	req *usrcreatesvc2.ConfirmRequest,
) (*usrcreatesvc2.ConfirmResponse, error) {

	userId, err := s.userCreation.ConfirmNewUser(ctx, req.Jwt)
	if err != nil {
		return nil, err
	}
	s.log.Info("user confirmed")

	return &usrcreatesvc2.ConfirmResponse{
		UserId: userId,
	}, nil
}
