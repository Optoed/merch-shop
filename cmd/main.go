package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"merch-shop/internal/handler"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/middleware"
	"merch-shop/pkg/config"
)

func main() {
	config.LoadEnv()

	database.InitDB()
	defer database.CloseDB()

	r := gin.Default()

	r.POST("/api/auth", handler.AuthHandler)

	protectedAPI := r.Group("/api")
	protectedAPI.Use(middleware.JWTMiddleware())
	{
		protectedAPI.GET("/info", handler.GetInfo)
		protectedAPI.POST("/sendCoin", handler.SendCoinHandler)
		protectedAPI.POST("/buy/:item", handler.BuyItem)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка при старте сервера:", err)
	}
}
