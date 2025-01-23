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

// TODO: connect grpc client for analytics microservice and sending email microservice here to that level

type Register interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
}

type serverAPI struct {
	regv1.UnimplementedRegistrationServer
	register Register
}

func RegisterServer(gRPC *grpc.Server, register Register) {
	regv1.RegisterRegistrationServer(gRPC, &serverAPI{register: register})
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *regv1.RegisterRequest,
) (*regv1.RegisterResponse, error) {

	// Validate request in regex
	if err := validateRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.register.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, register.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &regv1.RegisterResponse{
		UserId: userId,
	}, nil
}

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
