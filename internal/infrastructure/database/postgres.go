package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"merch-shop/pkg/config"
)

var DB *sqlx.DB

func InitDB(isTest bool) {
	DbUrl := config.GetDBUrl(isTest) // Используем флаг для тестовой базы
	var err error

	DB, err = sqlx.Connect("postgres", DbUrl)
	if err != nil {
		log.Fatalf("Error with connection to the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error with pinging the database: %v", err)
	}

	log.Println("Successful connection to the database!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection close")
	}
}

func ClearDB() {
	if DB != nil {
		// Учитывай порядок!
		tables := []string{
			"transactions",
			"inventory",
			"users",
		}

		for _, table := range tables {
			_, err := DB.Exec("DELETE FROM" + " " + table)
			if err != nil {
				log.Printf("Ошибка при очистке таблицы %s: %v", table, err)
			} else {
				log.Printf("Таблица %s успешно очищена", table)
			}
		}
	}
}
