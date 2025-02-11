package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CheckPasswordHash(password, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хешировании пароля: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}
