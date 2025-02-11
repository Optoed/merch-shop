package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"merch-shop/internal/infrastructure/config"
)

var DB *sqlx.DB

func InitDB() {
	DbUrl := config.GetDBUrl()

	conn, err := sqlx.Connect("postgres", DbUrl)
	if err != nil {
		log.Fatalf("Error with connection to the db: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("Error with pinging the db: %v", err)
	}

	log.Println("Successful connection to the db!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection close")
	}
}
