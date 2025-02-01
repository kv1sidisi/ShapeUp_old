package jwtsvc

import (
	"JWTService/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

type JWTService struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *JWTService {
	return &JWTService{log: log, cfg: cfg}
}

func (s *JWTService) GenerateAccessToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error) {
	operation = getOperationType(operation)
	if len(operation) == 0 {
		return "", errors.New("invalid operation type")
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
		return "", err
	}

	return accessToken, nil
}

func (s *JWTService) GenerateRefreshToken(ctx context.Context, uid int64, operation string, secretKey string) (string, error) {
	operation = getOperationType(operation)
	if len(operation) == 0 {
		return "", errors.New("invalid operation type")
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
		return "", err
	}

	return refreshToken, nil
}

func (s *JWTService) ValidateAccessToken(ctx context.Context, accessToken string, secretKey string) (uid int64, operation string, err error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, "", fmt.Errorf("token has expired")
			}
		} else {
			return 0, "", fmt.Errorf("expiration time not found in token")
		}

		uid = int64(claims["user_id"].(float64))
		operation = claims["operation"].(string)

		return uid, operation, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}

func (s *JWTService) ValidateRefreshToken(ctx context.Context, refreshToken string, secretKey string) (uid int64, operation string, err error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, "", fmt.Errorf("token has expired")
			}
		} else {
			return 0, "", fmt.Errorf("expiration time not found in token")
		}

		uid = int64(claims["user_id"].(float64))
		operation = claims["operation"].(string)

		return uid, operation, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}

func (s *JWTService) GenerateLink(ctx context.Context, linkBase string, uid int64, operation string, secretKey string) (string, error) {
	operation = getOperationType(operation)
	if len(operation) == 0 {
		return "", errors.New("invalid operation type")
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
