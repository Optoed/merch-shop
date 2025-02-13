package repository

import (
	"github.com/jmoiron/sqlx"
)

func UpdateBalanceTx(tx *sqlx.Tx, userID, amount int) error {
	query := `UPDATE users SET balance=balance+$1 WHERE id=$2`
	_, err := tx.Exec(query, amount, userID)
	return err
}

func CreateTransactionTx(tx *sqlx.Tx,
	senderID, receiverID int, receiverName string, amount int) error {
	query := `INSERT INTO transactions (sender_id, receiver_id, receiver_name, amount) 
			  VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(query, senderID, receiverID, receiverName, amount)
	return err
}
