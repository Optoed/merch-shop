package service

import (
	"database/sql"
	"errors"
	"merch-shop/internal/config"
	"merch-shop/internal/repository"
	"merch-shop/pkg/jwtAuth"
	"merch-shop/pkg/utils"
)

func Authenticate(username, password string) (string, error) {
	user, err := repository.GetUserByUsername(username)

	if errors.Is(err, sql.ErrNoRows) {
		return createUserAndGenerateJWT(username, password)
	} else if err != nil {
		return "", nil
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("Неправильный логин или пароль") //TODO
	}

	return generateJWT(user.ID)
}

func createUserAndGenerateJWT(username, password string) (string, error) {
	userID, err := repository.CreateUser(username, password, 1000)
	if err != nil {
		return "", err
	}
	return generateJWT(userID)
}

func generateJWT(userID int) (string, error) {
	token, err := jwtAuth.GenerateJWT(config.Cfg.SecretJWTKey, userID)
	if err != nil {
		return "", err
	}
	return token, nil
}
