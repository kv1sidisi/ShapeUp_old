package register

import (
	"RegistrationService/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Register struct {
	log       *slog.Logger
	userSaver UserSaver
	tokenTTL  time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

var (
	ErrUserExists = errors.New("user already exists")
)

// New returns a new instance of Register service.
func New(log *slog.Logger,
	userSaver UserSaver,
	tokenTTL time.Duration,
) *Register {
	return &Register{
		userSaver: userSaver,
		log:       log,
		tokenTTL:  tokenTTL,
	}
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (r *Register) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "register.RegisterNewUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := r.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			r.log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save new user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully saved and registered new user")

	return id, nil
}
