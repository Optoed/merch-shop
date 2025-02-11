package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"merch-shop/internal/config"
	"merch-shop/internal/handler"
	"merch-shop/internal/infrastructure/database"
)

func main() {
	config.LoadEnv()

	database.InitDB()
	defer database.CloseDB()

	r := gin.Default()

	r.POST("/api/auth", handler.AuthHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка при старте сервера:", err)
	}
}
