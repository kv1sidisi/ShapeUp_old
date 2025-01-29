package service

import (
	"context"
	"log/slog"
)

// AuthService struct represents the sending service and it is implementation of bottom layer of sending method of application.
type AuthService struct {
	log     *slog.Logger
	storage AuthManager
}

type AuthManager interface {
	FindUser(ctx context.Context,
		username string,
		password string,
	) (uid int64, err error)
}

// New returns a new instance of AuthService service.
func New(log *slog.Logger, storage AuthManager) *AuthService {
	return &AuthService{log: log, storage: storage}
}

func (as *AuthService) LoginUser(
	ctx context.Context,
	username string,
	password string,
) (userId int64, err error) {
	//TODO: try to find user in db
	//TODO: generate jwt and refresh tokens
	//TODO: return jwt and refresh tokens
	return 1, nil
}
