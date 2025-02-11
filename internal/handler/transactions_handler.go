package handler

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/models/requestModels"
	"net/http"
)

func SendCoinHandler(c *gin.Context) {
	senderName, exists := c.Get("username") //TODO пусть jwt кладет и userID, и username, и мы берем оба
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос."}) //TODO
		return
	}

	senderNameString, ok := senderName.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос."}) //TODO
		return
	}

	var (
		sendCoinRequest requestModels.SendCoinRequest
		err             error
	)
	if err = c.ShouldBindJSON(&sendCoinRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос."})
		return
	}

	// TODO допиши, тут вызываешь сервис и транзакции
}
