package handler

import (
	"github.com/gin-gonic/gin"
	"merch-shop/internal/models/requestModels"
	"merch-shop/internal/service"
	"net/http"
)

func AuthHandler(c *gin.Context) {
	var authRequest requestModels.AuthRequest
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()})
		return
	}

	token, err := service.Authenticate(authRequest.Username, authRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": "Успешная аутентификация.", "token": token})
}
