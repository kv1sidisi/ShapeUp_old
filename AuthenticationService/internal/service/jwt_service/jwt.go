package jwt_service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	accessTokenExpireTime  = time.Minute * 30
	refreshTokenExpireTime = time.Hour * 24 * 30
)

func GenerateAccessToken(userId int64, secretKey string) (accessToken string, err error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(accessTokenExpireTime).Unix(), // Expiration time
		"iat":     time.Now().Unix(),                            // Current time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId int64, secretKey string) (refreshToken string, err error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(refreshTokenExpireTime).Unix(), // Expiration time
		"iat":     time.Now().Unix(),                             // Current time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func ValidateAccessToken(accessToken string, secretKey string) (userId int64, err error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, fmt.Errorf("token has expired")
			}
		} else {
			return 0, fmt.Errorf("expiration time not found in token")
		}

		userId = int64(claims["user_id"].(float64))

		return userId, nil
	}

	return 0, fmt.Errorf("invalid token")
}

func ValidateRefreshToken(refreshToken string, secretKey string) (userId int64, err error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return 0, fmt.Errorf("token has expired")
			}
		} else {
			return 0, fmt.Errorf("expiration time not found in token")
		}

		userId = int64(claims["user_id"].(float64))

		return userId, nil
	}

	return 0, fmt.Errorf("invalid token")
}

func GenerateConfirmationLink(userId int64, secretKey string, linkBase string) (link string, err error) {
	return "", nil
}
