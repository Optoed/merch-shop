package main

import (
	"merch-shop/internal/infrastructure/config"
	"merch-shop/internal/infrastructure/db"
)

func main() {
	config.LoadEnv()

	db.InitDB()
	defer db.CloseDB()

}
