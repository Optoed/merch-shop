package handler

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/models/requestModels"
	"merch-shop/internal/service"
	"net/http"
)

func SendCoinHandler(c *gin.Context) {
	senderID, existsID := c.Get("user_id")
	if !existsID {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
		return
	}

	senderIDInt, okID := senderID.(int)
	if !okID {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
		return
	}

	senderName, existsName := c.Get("username")
	if !existsName {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Внутренняя ошибка сервера."})
	}

	var (
		sendCoinRequest requestModels.SendCoinRequest
		err             error
	)
	if err = c.ShouldBindJSON(&sendCoinRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос."})
		return
	}

	err = service.SendCoin(senderIDInt, senderName.(string), sendCoinRequest.ReceiverName, sendCoinRequest.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Успешный ответ."})
}
