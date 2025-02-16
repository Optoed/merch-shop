package unit

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"merch-shop/internal/handler"
	"merch-shop/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockBuyItem(userID int, item string) error {
	// Эмулируем успешную покупку
	if item == "invalidItem" {
		return fmt.Errorf("товар не найден")
	}
	return nil
}

func TestBuyItemHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем сервис BuyItem
	service.BuyItem = mockBuyItem

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)        // Устанавливаем user_id в контекст
	c.Set("username", "user1") // Устанавливаем username в контекст
	c.Params = append(c.Params, gin.Param{Key: "item", Value: "t-shirt"})

	requestBody := bytes.NewBufferString(`{"amount": 2}`)
	c.Request, _ = http.NewRequest("POST", "/buy/t-shirt", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.BuyItem(c)

	// Проверка результатов
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Успешный ответ.")
}

func TestBuyItemHandler_UserIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка запроса с ошибкой получения user_id из контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "user1") // Устанавливаем только username в контекст
	c.Params = append(c.Params, gin.Param{Key: "item", Value: "t-shirt"})

	requestBody := bytes.NewBufferString(`{"amount": 2}`)
	c.Request, _ = http.NewRequest("POST", "/buy/t-shirt", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.BuyItem(c)

	// Проверка результатов
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
}

func TestBuyItemHandler_ItemNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем ошибку в сервисе (товар не найден)
	service.BuyItem = func(userID int, item string) error {
		return fmt.Errorf("товар не найден")
	}

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)        // Устанавливаем user_id в контекст
	c.Set("username", "user1") // Устанавливаем username в контекст
	c.Params = append(c.Params, gin.Param{Key: "item", Value: "invalidItem"})

	requestBody := bytes.NewBufferString(`{"amount": 2}`)
	c.Request, _ = http.NewRequest("POST", "/buy/invalidItem", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.BuyItem(c)

	// Проверка результатов
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "товар не найден")
}
