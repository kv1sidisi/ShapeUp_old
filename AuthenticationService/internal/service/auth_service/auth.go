package auth_service

import (
	"AuthenticationService/api/pb/jwt_service"
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

// AuthService struct represents the sending service and it is implementation of bottom layer of sending method of application.
type AuthService struct {
	log     *slog.Logger
	cfg     *config.Config
	storage AuthManager
}

type AuthManager interface {
	FindUserByEmail(ctx context.Context,
		email string,
	) (user models.User, err error)
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
	jwtClient jwt_service.JWTClient,
) (userId int64, accessToken string, refreshToken string, err error) {
	//TODO: check user confirmed or not

	user, err := as.storage.FindUserByEmail(ctx, username)
	if err != nil {
		return 0, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return 0, "", "", fmt.Errorf("invalid credentials")
	}

	accessTokenGenResp, err := jwtClient.GenerateAccessToken(ctx, &jwt_service.AccessTokenRequest{
		Uid:       userId,
		Operation: accessOperationType,
	})
	if err != nil {
		return 0, "", "", err
	}
	accessToken = accessTokenGenResp.GetToken()

	refreshTokenGenResp, err := jwtClient.GenerateRefreshToken(ctx, &jwt_service.RefreshTokenRequest{
		Uid:       userId,
		Operation: refreshOperationType,
	})
	if err != nil {
		return 0, "", "", err
	}
	refreshToken = refreshTokenGenResp.GetToken()

	fmt.Println(userId, accessToken, refreshToken)

	if err := as.storage.AddSession(ctx, user.ID, accessToken, refreshToken); err != nil {
		return 0, "", "", err
	}

	return user.ID, accessToken, refreshToken, nil
}
