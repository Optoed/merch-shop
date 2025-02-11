package jwtAuth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"merch-shop/pkg/config"
	"time"
)

func GenerateJWT(secretKey string, userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method") //TODO Неавторизован.
		}
		return []byte(config.Cfg.SecretJWTKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		username := claims["username"].(string)
		return userID, username, nil
	}
	return 0, "", errors.New("Invalid token claims") //TODO Неавторизован.
}
