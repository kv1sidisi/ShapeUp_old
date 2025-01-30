package service

import (
	"AuthenticationService/internal/config"
	"AuthenticationService/internal/service/helpers"
	"context"
	"github.com/asaskevich/govalidator"
	"log/slog"
)

// AuthService struct represents the sending service and it is implementation of bottom layer of sending method of application.
type AuthService struct {
	log     *slog.Logger
	cfg     *config.Config
	storage AuthManager
}

type AuthManager interface {
	FindUserByName(ctx context.Context,
		username string,
		password string,
	) (uid int64, err error)
	FindUserByEmail(ctx context.Context,
		username string,
		password string,
	) (uid int64, err error)
	AddSession(ctx context.Context,
		uid int64,
		accessToken string,
		refreshToken string,
	) (err error)
}

// New returns a new instance of AuthService service.
func New(log *slog.Logger, cfg *config.Config, storage AuthManager) *AuthService {
	return &AuthService{log: log,
		storage: storage,
		cfg:     cfg}
}

func (as *AuthService) LoginUser(
	ctx context.Context,
	username string,
	password string,
) (userId int64, accessToken string, refreshToken string, err error) {
	// TODO: try pattern chain of responsibility
	if govalidator.IsEmail(username) {
		// if username is email
		userId, err = as.storage.FindUserByEmail(ctx, username, password)
	} else {
		// if username is login
		userId, err = as.storage.FindUserByName(ctx, username, password)
	}

	if err != nil {
		return 0, "", "", err
	}

	accessToken, err = jwt.GenerateAccessToken(userId, as.cfg.JWTSecret)
	if err != nil {
		return 0, "", "", err
	}
	refreshToken, err = jwt.GenerateRefreshToken(userId, as.cfg.JWTSecret)
	if err != nil {
		return 0, "", "", err
	}

	if err := as.storage.AddSession(ctx, userId, accessToken, refreshToken); err != nil {
		return 0, "", "", err
	}
	return userId, accessToken, refreshToken, nil
}
