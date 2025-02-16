package handler

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/service"
	"net/http"
)

func BuyItem(c *gin.Context) {
	item := c.Param("item")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос"})
		return
	}

	err := service.BuyItem(userID.(int), item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Успешный ответ."})
}
