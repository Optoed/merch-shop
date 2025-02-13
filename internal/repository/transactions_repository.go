package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/models"
)

func CreateTransactionTx(tx *sqlx.Tx,
	senderID, receiverID int, receiverName string, amount int) error {
	query := `INSERT INTO transactions (sender_id, receiver_id, receiver_name, amount) 
			  VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(query, senderID, receiverID, receiverName, amount)
	return err
}

func GetTransactionsFromUser(userID int) ([]models.Transaction, error) {
	var transactionsFromUser []models.Transaction
	query := `SELECT * FROM transactions WHERE sender_id=$1`

	err := database.DB.Select(&transactionsFromUser, query, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return transactionsFromUser, nil
}

func GetTransactionsToUser(userID int) ([]models.Transaction, error) {
	var transactionsToUser []models.Transaction
	query := `SELECT * FROM transactions WHERE receiver_id=$1`

	err := database.DB.Select(&transactionsToUser, query, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return transactionsToUser, nil
}
