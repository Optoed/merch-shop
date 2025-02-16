package handler

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/repository"
	"net/http"
)

func GetInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неавторизован."})
		return
	}

	userIdInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неавторизован."})
		return
	}

	// 0) balance
	balance, err := repository.GetUserBalanceByID(userIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
		return
	}

	// 1) inventory
	inventory, err := repository.GetUserInventory(userIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
		return
	}

	// 2) transactions to current user
	transactionToUser, err := repository.GetTransactionsToUser(userIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
		return
	}

	// 3) transactions from current user
	transactionFromUser, err := repository.GetTransactionsFromUser(userIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"description": "Успешный ответ.",
		"schema": gin.H{
			"coins":     balance,
			"inventory": inventory,
			"coinHistory": gin.H{
				"received": transactionFromUser,
				"sent":     transactionToUser,
			},
		},
	})
}
