package authsvc

import (
	pbjwtsvc "AuthenticationService/api/pb/jwtsvc"
	pbsendsvc "AuthenticationService/api/pb/sendsvc"
	"AuthenticationService/internal/config"
	"AuthenticationService/internal/domain/models"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

const (
	refreshOperationType = "refresh"
	accessOperationType  = "access"
)

// AuthSvc struct represents the sending service and it is implementation of bottom layer of sending method of application.
type AuthSvc struct {
	log           *slog.Logger
	cfg           *config.Config
	storage       AuthMgr
	sendingClient pbsendsvc.SendingClient
	jwtClient     pbjwtsvc.JWTClient
}

type AuthMgr interface {
	FindUserByEmail(ctx context.Context,
		email string,
	) (user models.User, err error)
	AddSession(ctx context.Context,
		uid int64,
		accessToken string,
		refreshToken string,
	) (err error)
	IsUserConfirmed(ctx context.Context,
		uid int64,
	) (confirmed bool, err error)
}

// New returns a new instance of AuthSvc service.
func New(log *slog.Logger,
	cfg *config.Config,
	storage AuthMgr,
	sendingClient pbsendsvc.SendingClient,
	jwtClient pbjwtsvc.JWTClient) *AuthSvc {
	return &AuthSvc{log: log,
		storage:       storage,
		cfg:           cfg,
		sendingClient: sendingClient,
		jwtClient:     jwtClient}

}

func (as *AuthSvc) LoginUser(
	ctx context.Context,
	username string,
	password string,
) (userId int64, accessToken string, refreshToken string, err error) {
	const op = "authsvc.LoginUser"
	log := as.log.With(slog.String("op", op))

	user, err := as.storage.FindUserByEmail(ctx, username)
	if err != nil {
		return 0, "", "", err
	}

	isConfirmed, err := as.storage.IsUserConfirmed(ctx, user.ID)
	if err != nil {
		return 0, "", "", err
	}
	if !isConfirmed {
		return 0, "", "", fmt.Errorf("user is not confirmed")
	}

	log.Info("user is confirmed")

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return 0, "", "", fmt.Errorf("invalid credentials")
	}

	accessTokenGenResp, err := as.jwtClient.GenerateAccessToken(ctx, &pbjwtsvc.AccessTokenRequest{
		Uid:       user.ID,
		Operation: accessOperationType,
	})
	if err != nil {
		return 0, "", "", err
	}
	accessToken = accessTokenGenResp.GetToken()
	log.Info("access token generated", slog.String("accessToken", accessToken))

	refreshTokenGenResp, err := as.jwtClient.GenerateRefreshToken(ctx, &pbjwtsvc.RefreshTokenRequest{
		Uid:       user.ID,
		Operation: refreshOperationType,
	})
	if err != nil {
		return 0, "", "", err
	}
	refreshToken = refreshTokenGenResp.GetToken()
	log.Info("refresh token generated", slog.String("refreshToken", refreshToken))

	if err := as.storage.AddSession(ctx, user.ID, accessToken, refreshToken); err != nil {
		return 0, "", "", err
	}
	log.Info("session created")

	return user.ID, accessToken, refreshToken, nil
}
