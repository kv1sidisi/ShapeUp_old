package register

import (
	regv1 "RegistrationService/api/pb"
	"RegistrationService/internal/config"
	"RegistrationService/internal/storage"
	"RegistrationService/pkg/utils/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// UserCreation interface represents upper layer of userCreation method of application.
type UserCreation interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)

	ConfirmNewUser(
		ctx context.Context,
		userId int64,
	) (err error)
}

// serverAPI represents the handler for the gRPC server.
type serverAPI struct {
	regv1.UnimplementedUserCreationServer
	userCreation UserCreation
	cfg          *config.Config
}

// RegisterServer registers the request handler for registration in the gRPC server.
func RegisterServer(gRPC *grpc.Server, userCreation UserCreation, cfg *config.Config) {
	regv1.RegisterUserCreationServer(gRPC, &serverAPI{userCreation: userCreation, cfg: cfg})
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

	// UserCreation the new user
	userId, err := s.userCreation.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	// Generate link for account confirmation
	link, err := jwt.JwtLinkGeneration(userId, s.cfg.JWTSecret)

	fmt.Printf(link)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	// TODO: send link with email sender service

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

// ConfirmAccount is the gRPC server handler method, the top layer of the registration process.
func (s *serverAPI) ConfirmAccount(ctx context.Context,
	req *regv1.ConfirmRequest,
) (*regv1.ConfirmResponse, error) {

	userId, err := jwt.VerifyToken(req.Jwt, s.cfg.JWTSecret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if err := s.userCreation.ConfirmNewUser(ctx, userId); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &regv1.ConfirmResponse{
		UserId: userId,
	}, nil
}
