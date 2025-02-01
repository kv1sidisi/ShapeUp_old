package user_creation

import (
	"RegistrationService/api/pb/jwt_service"
	"RegistrationService/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// UserCreation struct represents the registration service and it is implementation of bottom layer of register method of application.
type UserCreation struct {
	log       *slog.Logger
	userSaver UserManager
}

// UserManager interface defines the method for saving user information in database.
type UserManager interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
	ConfirmAccount(
		ctx context.Context,
		uid int64,
	) (err error)
}

// New returns a new instance of UserCreation service.
func New(log *slog.Logger,
	userSaver UserManager,
) *UserCreation {
	return &UserCreation{
		userSaver: userSaver,
		log:       log,
	}
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (r *UserCreation) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "register.RegisterNewUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	// Generate a hashed password from the provided password.
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("password hash generated")

	log.Info("saving new user to database")
	id, err := r.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			r.log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed to save new user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// ConfirmNewUser confirms account
// If user does not exist returns error
func (r *UserCreation) ConfirmNewUser(ctx context.Context, token string, jwtClient jwt_service.JWTClient) (userId int64, err error) {
	const op = "register.ConfirmAccount"

	log := r.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)

	log.Info("verifying confirmation token")
	validationResp, err := jwtClient.ValidateAccessToken(ctx, &jwt_service.ValidateAccessTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Error("failed to verify confirmation token", err)
		return -1, status.Error(codes.Unauthenticated, "invalid token")
	}
	userId = validationResp.GetUid()

	if err := r.userSaver.ConfirmAccount(ctx, userId); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			r.log.Warn("user does not exist")
			return -1, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed to confirm new user")
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
