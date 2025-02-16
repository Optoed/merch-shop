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

func mockSendCoin(senderID int, senderName, receiverName string, amount int) error {
	if receiverName == "invalidUser" {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}

func TestSendCoinHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка мока для функции SendCoin
	service.SendCoin = mockSendCoin

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)        // Устанавливаем user_id в контекст
	c.Set("username", "user1") // Устанавливаем username в контекст

	requestBody := bytes.NewBufferString(`{"toUser": "user2", "amount": 100}`)
	c.Request, _ = http.NewRequest("POST", "/sendCoin", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.SendCoinHandler(c)

	// Проверка результатов
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Успешный ответ.")
}

func TestSendCoinHandler_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка запроса с некорректным телом
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)        // Устанавливаем user_id в контекст
	c.Set("username", "user1") // Устанавливаем username в контекст

	requestBody := bytes.NewBufferString(`{"toUser": "user2"}`)
	c.Request, _ = http.NewRequest("POST", "/sendCoin", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.SendCoinHandler(c)

	// Проверка результатов
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос.")
}

func TestSendCoinHandler_UserIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка запроса с ошибкой получения user_id из контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "user1") // Устанавливаем только username в контекст

	requestBody := bytes.NewBufferString(`{"toUser": "user2", "amount": 100}`)
	c.Request, _ = http.NewRequest("POST", "/sendCoin", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.SendCoinHandler(c)

	// Проверка результатов
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Внутренняя ошибка сервера.")
}

func TestSendCoinHandler_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем ошибку в сервисе (не найден пользователь)
	service.SendCoin = func(senderID int, senderName, receiverName string, amount int) error {
		return fmt.Errorf("пользователь не найден")
	}

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)        // Устанавливаем user_id в контекст
	c.Set("username", "user1") // Устанавливаем username в контекст

	requestBody := bytes.NewBufferString(`{"toUser": "invalidUser", "amount": 100}`)
	c.Request, _ = http.NewRequest("POST", "/sendCoin", requestBody)
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов обработчика
	handler.SendCoinHandler(c)

	// Проверка результатов
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "пользователь не найден")
}
