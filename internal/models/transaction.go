package models

type Transaction struct {
	ID           int    `json:"id" db:"id"`
	SenderID     int    `json:"sender_id" db:"sender_id"`
	SenderName   string `json:"sender_name" db:"sender_name"`
	ReceiverID   int    `json:"receiver_id" db:"receiver_id"`
	ReceiverName string `json:"receiver_name" db:"receiver_name"`
	Amount       int    `json:"amount" db:"amount"`
}

type TransactionTo struct {
	ReceiverName string `json:"toUser" db:"receiver_name"`
	Amount       int    `json:"amount" db:"amount"`
}

type TransactionFrom struct {
	SenderName string `json:"fromUser" db:"sender_name"`
	Amount     int    `json:"amount" db:"amount"`
}
