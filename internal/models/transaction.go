package models

import "time"

type Transaction struct {
	ID           int       `json:"id" db:"id"`
	SenderID     int       `json:"sender_id" db:"sender_id"`
	ReceiverID   int       `json:"receiver_id" db:"receiver_id"`
	ReceiverName string    `json:"receiver_name" db:"receiver_name"`
	Amount       int       `json:"amount" db:"amount"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
}
