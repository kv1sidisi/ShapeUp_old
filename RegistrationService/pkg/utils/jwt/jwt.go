package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

const (
	operation_type = "confirmation"
)

// JwtLinkGeneration generates link from jwt_service token to confirm account.
func JwtLinkGeneration(userId int64, secretKey string) (string, error) {
	token, err := generateToken(userId, secretKey)
	if err != nil {
		return "", err
	}

	link := "http://localhost:8082/confirm_account?token=" + token
	return link, nil
}

// VerifyToken verifies jwt_service token and returns user ID.
func VerifyToken(log *slog.Logger, tokenString string, secretKey string) (int64, error) {
	op := "jwt_service.VerifyToken"

	log.With(slog.String("op", op))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error("token is not parsed", tokenString)
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	log.Info("parsed token")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, fmt.Errorf("token has expired")
			}
		} else {
			return 0, fmt.Errorf("expiration time not found in token")
		}

		log.Info("token not expired")

		if claims["operation"].(string) != operation_type {
			return -1, fmt.Errorf("invalid operation type")
		}

		log.Info("token operation type is correct")

		userId := int64(claims["user_id"].(float64))

		log.Info("successfully verified, userId is ", userId)
		return userId, nil
	}

	return -1, fmt.Errorf("invalid token")
}

// generateToken generates JWT.
func generateToken(userId int64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userId,
		"operation": operation_type,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Expiration time
		"iat":       time.Now().Unix(),                     // Current time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signing token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TODO: support token generation for different operations
func verifyOperationType() {}
