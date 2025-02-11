package repository

import (
	"log"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/models"
	"merch-shop/pkg/utils"
)

func CreateUser(username, password string, balance int) (int, error) {
	query := `INSERT INTO users (username, password_hash, balance)
			  VALUES ($1, $2, $3) RETURNING id`
	var userID int

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return -1, err
	}

	err = database.DB.QueryRow(query, username, passwordHash, balance).Scan(&userID)

	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return -1, err
	}
	log.Println("Пользователь успешно создан")
	return userID, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT * FROM users WHERE username=$1`
	var user models.User
	err := database.DB.Get(&user, query, username)
	//log.Println("user GetUserByUsername = ", user, "error = ", err)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
