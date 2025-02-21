package authsvc

import (
	"context"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	pbjwtsvc "github.com/kv1sidisi/shapeup/pkg/proto/jwtsvc/pb"
	pbsendsvc "github.com/kv1sidisi/shapeup/pkg/proto/sendsvc/pb"
	"github.com/kv1sidisi/shapeup/services/authsvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/authsvc/cmd/grpccl/consts"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/config"
	"github.com/kv1sidisi/shapeup/services/authsvc/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

const (
	refreshOperationType = "refresh"
	accessOperationType  = "access"
)

// AuthSvc authentication service.
type AuthSvc struct {
	log           *slog.Logger
	cfg           *config.Config
	storage       AuthMgr
	sendingClient pbsendsvc.SendingClient
	jwtClient     pbjwtsvc.JWTClient
}

// AuthMgr manager for database.
type AuthMgr interface {
	FindUserByEmail(ctx context.Context,
		email string,
	) (user models.User, err error)
	AddSession(ctx context.Context,
		uid []byte,
		refreshToken string,
	) (err error)
	IsUserConfirmed(ctx context.Context,
		uid []byte,
	) (confirmed bool, err error)
}

func New(log *slog.Logger,
	cfg *config.Config,
	storage AuthMgr,
	grpccl *grpccl.GRPCClients,
) *AuthSvc {
	return &AuthSvc{log: log,
		storage:       storage,
		cfg:           cfg,
		sendingClient: grpccl.Cl[consts.SendSvc].Client.(pbsendsvc.SendingClient),
		jwtClient:     grpccl.Cl[consts.JWTSvc].Client.(pbjwtsvc.JWTClient),
	}

}

// LoginUser logs user in.
//
// Returns:
//
//   - userId, access and refresh tokens if successful.
//
//   - Error if: Database fails.
//     User is not confirmed.
//     Invalid password.
//     Failed to generate JWT access or refresh tokens through JWTService
//     Failed to add new session for user.
func (as *AuthSvc) LoginUser(
	ctx context.Context,
	username string,
	password string,
) (uif []byte, accessToken string, refreshToken string, err error) {
	const op = "authsvc.LoginUser"
	log := as.log.With(slog.String("op", op),
		slog.String("username", username))

	user, err := as.storage.FindUserByEmail(ctx, username)
	if err != nil {
		return nil, "", "", err
	}

	isConfirmed, err := as.storage.IsUserConfirmed(ctx, user.ID)
	if err != nil {
		return nil, "", "", err
	}
	if !isConfirmed {
		log.Error("user is not confirmed")
		return nil, "", "", errdefs.ErrUserNotConfirmed
	}

	log.Info("user is confirmed")

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid password")
		return nil, "", "", errdefs.InvalidCredentials
	}

	accessTokenGenResp, err := as.jwtClient.GenerateToken(ctx, &pbjwtsvc.GenerateTokenRequest{
		Uid:       user.ID,
		Operation: accessOperationType,
	})
	if err != nil {
		return nil, "", "", err
	}
	accessToken = accessTokenGenResp.GetToken()
	log.Info("access token generated", slog.String("accessToken", accessToken))

	refreshTokenGenResp, err := as.jwtClient.GenerateToken(ctx, &pbjwtsvc.GenerateTokenRequest{
		Uid:       user.ID,
		Operation: refreshOperationType,
	})
	if err != nil {
		return nil, "", "", err
	}
	refreshToken = refreshTokenGenResp.GetToken()
	log.Info("refresh token generated", slog.String("refreshToken", refreshToken))

	if err := as.storage.AddSession(ctx, user.ID, refreshToken); err != nil {
		return nil, "", "", err
	}
	log.Info("session created")

	return user.ID, accessToken, refreshToken, nil
}
