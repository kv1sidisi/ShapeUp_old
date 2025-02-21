package grpcsrv

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	usrcreatesvc "github.com/kv1sidisi/shapeup/pkg/proto/usercreatesvc/pb"
	"github.com/kv1sidisi/shapeup/services/regsvc/internal/config"
	"google.golang.org/grpc"
	"log/slog"
	"strings"
)

// UsrCreateSvc service for serverAPI.
type UsrCreateSvc interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (uid []byte, err error)

	ConfirmNewUser(
		ctx context.Context,
		jwt string,
	) (uid []byte, err error)
}

// serverAPI handler for the gRPC server.
type serverAPI struct {
	usrcreatesvc.UnimplementedUserCreationServer
	userCreation UsrCreateSvc
	cfg          *config.Config
	log          *slog.Logger
}

// RegisterServer registers services in the GRPC server.
//
// Returns serverAPI as handler for GRPC server.
func RegisterServer(gRPC *grpc.Server,
	userCreation UsrCreateSvc,
	cfg *config.Config,
	log *slog.Logger,
) {
	usrcreatesvc.RegisterUserCreationServer(
		gRPC,
		&serverAPI{
			userCreation: userCreation,
			cfg:          cfg,
			log:          log,
		})
}

// Register is the GRPC server handler method. Registers new user.
//
// Returns:
//
//   - A pointer to RegisterResponse if successful.
//
//   - An error if: Request is invalid. Error while registering user through service.
func (s *serverAPI) Register(
	ctx context.Context,
	req *usrcreatesvc.RegisterRequest,
) (*usrcreatesvc.RegisterResponse, error) {
	const op = "grpcsrv.Register"
	log := s.log.With(slog.String("op", op))

	// Validate request in regex
	if err := validateRegisterRequest(log, req); err != nil {
		return nil, err
	}
	log.Info("register request valid")

	// UsrCreateSvc the new user
	uid, err := s.userCreation.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	log.Info("user registered successfully")

	// Return the response with the user ID
	return &usrcreatesvc.RegisterResponse{
		Uid: uid,
	}, nil
}

// validateRegisterRequest performs validation on the registration request.
//
// Returns nil if request is invalid.
func validateRegisterRequest(log *slog.Logger, req *usrcreatesvc.RegisterRequest) error {
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

// Confirm is the gRPC server handler method. Confirms user account.
// Returns:
//
//   - A pointer to ConfirmResponse if successful.
//
//   - An error if: Request is invalid. Error while confirming user through service.
func (s *serverAPI) Confirm(ctx context.Context,
	req *usrcreatesvc.ConfirmRequest,
) (*usrcreatesvc.ConfirmResponse, error) {
	const op = "grpcsrv.Confirm"
	log := s.log.With(slog.String("op", op))

	if len(req.Jwt) == 0 {
		return nil, errdefs.InvalidCredentials
	}

	uid, err := s.userCreation.ConfirmNewUser(ctx, req.Jwt)
	if err != nil {
		return nil, err
	}
	log.Info("user confirmed")

	return &usrcreatesvc.ConfirmResponse{
		Uid: uid,
	}, nil
}
