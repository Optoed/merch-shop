package models

type Item struct {
	Name string `json:"name" db:"name"`
	Cost int    `json:"cost" db:"cost"`
}
