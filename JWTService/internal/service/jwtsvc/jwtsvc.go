package jwtsvc

import (
	"JWTService/internal/config"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

const (
	refreshOperationType      = "refresh"
	accessOperationType       = "access"
	confirmationOperationType = "confirmation"

	accessTokenExpireTime  = time.Minute * 30
	refreshTokenExpireTime = time.Hour * 24 * 30
)

const (
	InvalidOperationType = "invalid operation type"
	InvalidSigningMethod = "invalid token signing method"
	TokenExpired         = "token expired"
	InvalidToken         = "invalid token"
)

type JWTSvc struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *JWTSvc {
	return &JWTSvc{log: log, cfg: cfg}
}

func (s *JWTSvc) GenerateAccessToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error) {
	operation = getOperationType(operation)
	if len(operation) == 0 {
		return "", status.Error(codes.InvalidArgument, InvalidOperationType)
	}

	claims := jwt.MapClaims{
		"user_id":   uid,
		"operation": operation,
		"exp":       time.Now().Add(accessTokenExpireTime).Unix(), // Expiration time
		"iat":       time.Now().Unix(),                            // Current time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return accessToken, nil
}

func (s *JWTSvc) GenerateRefreshToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error) {
	operation = getOperationType(operation)
	if len(operation) == 0 {
		return "", status.Error(codes.InvalidArgument, InvalidOperationType)
	}

	claims := jwt.MapClaims{
		"user_id":   uid,
		"operation": operation,
		"exp":       time.Now().Add(refreshTokenExpireTime).Unix(), // Expiration time
		"iat":       time.Now().Unix(),                             // Current time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return refreshToken, nil
}

func (s *JWTSvc) ValidateAccessToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, InvalidSigningMethod)
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, "", status.Error(codes.Unauthenticated, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, "", status.Error(codes.InvalidArgument, TokenExpired)
			}
		} else {
			return 0, "", status.Error(codes.InvalidArgument, TokenExpired)
		}

		uid = int64(claims["user_id"].(float64))
		operation = claims["operation"].(string)

		return uid, operation, nil
	}

	return 0, "", status.Error(codes.InvalidArgument, InvalidToken)
}

func (s *JWTSvc) ValidateRefreshToken(ctx context.Context, refreshToken string, secretKey string) (uid int64, operation string, err error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, InvalidSigningMethod)
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, "", status.Error(codes.Unauthenticated, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, "", status.Error(codes.InvalidArgument, TokenExpired)
			}
		} else {
			return 0, "", status.Error(codes.InvalidArgument, TokenExpired)
		}

		uid = int64(claims["user_id"].(float64))
		operation = claims["operation"].(string)

		return uid, operation, nil
	}

	return 0, "", status.Error(codes.InvalidArgument, InvalidToken)
}

func (s *JWTSvc) GenerateLink(ctx context.Context, linkBase string, uid int64, operation string, secretKey string) (string, error) {
	operation = getOperationType(operation)
	if len(operation) == 0 {
		return "", status.Error(codes.InvalidArgument, InvalidOperationType)
	}

	token, err := s.GenerateAccessToken(ctx, uid, operation, secretKey)
	if err != nil {
		return "", err
	}

	link := linkBase + token

	return link, nil
}

func getOperationType(operation string) string {
	switch operation {
	case refreshOperationType:
		return refreshOperationType
	case accessOperationType:
		return accessOperationType
	case confirmationOperationType:
		return confirmationOperationType
	}

	return ""
}
