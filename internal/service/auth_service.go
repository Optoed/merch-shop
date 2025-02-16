package service

import (
	"database/sql"
	"errors"
	"merch-shop/internal/repository"
	"merch-shop/pkg/config"
	"merch-shop/pkg/jwtAuth"
	"merch-shop/pkg/utils"
)

func Authenticate(username, password string) (string, error) {
	user, err := repository.GetUserByUsername(username)

	if errors.Is(err, sql.ErrNoRows) {
		return createUserAndGenerateJWT(username, password)
	} else if err != nil {
		return "", errors.New("Неавторизован.")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("Неавторизован.")
	}

	return generateJWT(user.ID, username)
}

func createUserAndGenerateJWT(username, password string) (string, error) {
	userID, err := repository.CreateUser(username, password, 1000)
	if err != nil {
		return "", errors.New("внутренняя ошибка сервера")
	}
	return generateJWT(userID, username)
}

func generateJWT(userID int, username string) (string, error) {
	token, err := jwtAuth.GenerateJWT(config.Cfg.SecretJWTKey, userID, username)
	if err != nil {
		return "", errors.New("внутренняя ошибка сервера")
	}
	return token, nil
}
