package register

import (
	regv1 "RegistrationService/api/pb"
	"RegistrationService/internal/config"
	"RegistrationService/internal/service/register"
	"RegistrationService/pkg/utils/jwt"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// Register interface represents upper layer of register method of application.
type Register interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	regv1.UnimplementedRegistrationServer
	register Register
	cfg      *config.Config
}

// RegisterServer registers the request handler for registration in the gRPC server.
func RegisterServer(gRPC *grpc.Server, register Register, cfg *config.Config) {
	regv1.RegisterRegistrationServer(gRPC, &serverAPI{register: register, cfg: cfg})
}

// Register is the gRPC server handler method, the top layer of the registration process.
func (s *serverAPI) Register(
	ctx context.Context,
	req *regv1.RegisterRequest,
) (*regv1.RegisterResponse, error) {

	// Validate request in regex
	if err := validateRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Register the new user
	userId, err := s.register.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, register.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	// JWT generation
	jwtToken, err := jwt.GenerateToken(userId, s.cfg.JWTSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "jwt generation error")
	}

	// Confirmation link generation with JWT

	// Return the response with the user ID
	return &regv1.RegisterResponse{
		UserId: userId,
	}, nil
}

// validateRegisterRequest performs validation on the registration request.
func validateRegisterRequest(req *regv1.RegisterRequest) error {
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
