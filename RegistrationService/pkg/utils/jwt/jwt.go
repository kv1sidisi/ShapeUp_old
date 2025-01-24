package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	operation_type = "confirmation"
)

// GenerateToken generates JWT for account confirmation.
func GenerateToken(userId int64, secretKey string) (string, error) {
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

func VerifyToken(tokenString, secretKey string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["operation"] == operation_type {
			return -1, fmt.Errorf("invalid operation type")
		}

		userId := int64(claims["user_id"].(float64))
		return userId, nil
	}

	return -1, fmt.Errorf("invalid token")
}
