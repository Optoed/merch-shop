package repository

import (
	"log"
	"merch-shop/internal/infrastructure/db"
)

func СreateUser(username, passwordHash string, balance int) (int, error) {
	query := `INSERT INTO users (username, password_hash, balance)
			  VALUES (:username, :password_hash, :balance) RETURNING id`

	var userID int
	err := db.DB.QueryRow(query, username, passwordHash, balance).Scan(&userID)

	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return -1, err
	}

	log.Println("Пользователь успешно создан")
	return userID, nil
}
