package main

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/infrastructure/config"
	"merch-shop/internal/infrastructure/db"
)

func main() {
	config.LoadEnv()

	db.InitDB()
	defer db.CloseDB()

	r := gin.Default()
	r.Run("0.0.0.0:8080")
}
