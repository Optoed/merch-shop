package models

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Balance      int       `json:"balance" db:"balance"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
