package repository

import (
	"github.com/jmoiron/sqlx"
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

func GetUserByID(userID int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id=$1`
	var user models.User
	err := database.DB.Get(&user, query, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserBalanceByID(userID int) (int, error) {
	query := `SELECT balance FROM users WHERE id=$1`
	var balance int
	err := database.DB.Get(&balance, query, userID)
	if err != nil {
		return 0, nil
	}
	return balance, nil
}

func IncreaseBalanceByAmountTx(tx *sqlx.Tx, userID, amount int) error {
	query := `UPDATE users SET balance=balance+$1 WHERE id=$2`
	_, err := tx.Exec(query, amount, userID)
	return err
}

func DecreaseBalanceByAmountTx(tx *sqlx.Tx, userID, amount int) error {
	query := `UPDATE users SET balance=balance-$1 WHERE id=$2`
	_, err := tx.Exec(query, amount, userID)
	return err
}

func SetUserBalance(userID, amount int) error {
	query := `UPDATE users SET balance=$1 WHERE id=$2`
	_, err := database.DB.Exec(query, amount, userID)
	return err
}
