package handler

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/models/requestModels"
	"merch-shop/internal/service"
	"net/http"
)

func SendCoinHandler(c *gin.Context) {
	//senderName, existsName := c.Get("username") //TODO пусть jwt кладет и userID, и username, и мы берем оба
	senderID, existsID := c.Get("user_id")
	if !existsID { // !existsName ||
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос."}) //TODO
		return
	}

	//senderNameString, okName := senderName.(string)
	senderIDInt, okID := senderID.(int)
	if !okID { //!okName ||
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

	err = service.SendCoin(senderIDInt, sendCoinRequest.ReceiverName, sendCoinRequest.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()}) //TODO
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Успешный ответ."}) //TODO
}
