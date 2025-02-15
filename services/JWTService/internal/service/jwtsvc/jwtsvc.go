package jwtsvc

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/services/jwtsvc/internal/config"
	opinf "github.com/kv1sidisi/shapeup/services/jwtsvc/internal/service/operations"
	"log/slog"
	"time"
)

type JWTSvc struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *JWTSvc {
	return &JWTSvc{log: log, cfg: cfg}
}

func (s *JWTSvc) GenerateToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error) {
	const op = "jwtsvc.GenerateAccessToken"
	log := s.log.With(slog.String("op", op))

	opInfo := opinf.GetOperationInfo(operation)
	if len(operation) == 0 {
		log.Error("invalid operation type")
		return "", errdefs.InvalidOperationType
	}

	claims := jwt.MapClaims{
		"user_id":   uid,
		"operation": operation,
		"exp":       time.Now().Add(opInfo.ExpireTime).Unix(), // Expiration time
		"iat":       time.Now().Unix(),                        // Current time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Error("failed to sign access token")
		return "", errdefs.ErrInternal
	}

	return accessToken, nil
}

func (s *JWTSvc) ValidateToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error) {
	const op = "jwtsvc.ValidateAccessToken"
	log := s.log.With(slog.String("op", op))

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("invalid signing method")
			return nil, errdefs.InvalidSigningMethod
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error("failed to parse access token")
		return 0, "", errdefs.InvalidCredentials
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				log.Error("access token is expired")
				return 0, "", errdefs.ErrTokenExpired
			}
		} else {
			log.Error("error getting \"expired\" claim")
			return 0, "", errdefs.ErrInternal
		}

		uid = int64(claims["user_id"].(float64))
		operation = claims["operation"].(string)

		return uid, operation, nil
	}

	log.Error("invalid access token")
	return 0, "", errdefs.InvalidToken
}

func (s *JWTSvc) GenerateLink(ctx context.Context, linkBase string, uid int64, operation string, secretKey string) (string, error) {
	const op = "jwtsvc.GenerateLink"
	log := s.log.With(slog.String("op", op))

	token, err := s.GenerateToken(ctx, uid, operation, secretKey)
	if err != nil {
		return "", err
	}

	link := linkBase + token
	log.Info("generated", slog.String("link", link))

	return link, nil
}
