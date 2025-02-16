package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/models"
)

func CreateTransactionTx(
	tx *sqlx.Tx,
	senderID int,
	senderName string,
	receiverID int,
	receiverName string,
	amount int) error {

	query := `INSERT INTO transactions (sender_id, sender_name, receiver_id, receiver_name, amount) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(query, senderID, senderName, receiverID, receiverName, amount)
	return err
}

func GetTransactionsFromUser(userID int) ([]models.TransactionFrom, error) {
	var transactionsFromUser []models.TransactionFrom
	query := `SELECT sender_name, amount FROM transactions WHERE receiver_id=$1`

	err := database.DB.Select(&transactionsFromUser, query, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return transactionsFromUser, nil
}

func GetTransactionsToUser(userID int) ([]models.TransactionTo, error) {
	var transactionsToUser []models.TransactionTo
	query := `SELECT receiver_name, amount FROM transactions WHERE sender_id=$1`

	err := database.DB.Select(&transactionsToUser, query, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return transactionsToUser, nil
}
