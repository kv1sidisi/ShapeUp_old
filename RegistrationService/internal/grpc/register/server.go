package register

import (
	regv1 "RegistrationService/api/pb"
	"RegistrationService/internal/service/register"
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

//ATTENTION ATTENTION ATTENTION ATTENTION ATTENTION ATTENTION ATTENTION ATTENTION
// connect grpc client for analytics microservice and sending email microservice here to that level

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
}

// RegisterServer registers the request handler for registration in the gRPC server.
func RegisterServer(gRPC *grpc.Server, register Register) {
	regv1.RegisterRegistrationServer(gRPC, &serverAPI{register: register})
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
