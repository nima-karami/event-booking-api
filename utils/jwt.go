package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, userID int64, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	secretKey := getSecretKey()
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (userID int64, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		secretKey := getSecretKey()
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	floatUserID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("invalid user_id in token")
	}
	userID = int64(floatUserID)

	role, ok = claims["role"].(string)
	if !ok {
		return 0, "", errors.New("invalid role in token")
	}

	return userID, role, nil
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	return secret
}
