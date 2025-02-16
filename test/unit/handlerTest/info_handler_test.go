package unit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"merch-shop/internal/handler"
	"merch-shop/internal/models"
	"merch-shop/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockGetUserBalanceByID(userID int) (int, error) {
	// Эмулируем успешный ответ
	if userID == 1 {
		return 1000, nil
	}
	return 0, fmt.Errorf("баланс не найден")
}

func mockGetUserInventory(userID int) ([]models.InventoryItemResponse, error) {
	// Эмулируем успешный ответ с инвентарем
	if userID == 1 {
		return []models.InventoryItemResponse{
			{Type: "t-shirt", Quantity: 2},
			{Type: "hat", Quantity: 3},
		}, nil
	}
	return nil, fmt.Errorf("инвентарь не найден")
}

func mockGetTransactionsFromUser(userID int) ([]models.TransactionFrom, error) {
	// Эмулируем успешный ответ с транзакциями
	if userID == 1 {
		return []models.TransactionFrom{
			{SenderName: "user2", Amount: 50},
			{SenderName: "user3", Amount: 30},
		}, nil
	}
	return nil, fmt.Errorf("транзакции не найдены")
}

func mockGetTransactionsToUser(userID int) ([]models.TransactionTo, error) {
	// Эмулируем успешный ответ с транзакциями
	if userID == 1 {
		return []models.TransactionTo{
			{ReceiverName: "user2", Amount: 20},
			{ReceiverName: "user3", Amount: 10},
		}, nil
	}
	return nil, fmt.Errorf("транзакции не найдены")
}

func TestGetInfoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем все репозитории
	repository.GetUserBalanceByID = mockGetUserBalanceByID
	repository.GetUserInventory = mockGetUserInventory
	repository.GetTransactionsToUser = mockGetTransactionsToUser
	repository.GetTransactionsFromUser = mockGetTransactionsFromUser

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)        // Устанавливаем user_id в контекст
	c.Set("username", "user1") // Устанавливаем username в контекст

	// Вызов обработчика
	handler.GetInfo(c)

	// Ожидаемый результат в формате JSON
	expectedResponse := `{
		"description": "Успешный ответ.",
		"schema": {
			"coins": 1000,
			"inventory": [
				{"type": "t-shirt", "quantity": 2},
				{"type": "hat", "quantity": 3}
			],
			"coinHistory": {
				"received": [
					{"fromUser": "user2", "amount": 50},
					{"fromUser": "user3", "amount": 30}
				],
				"sent": [
					{"toUser": "user2", "amount": 20},
					{"toUser": "user3", "amount": 10}
				]
			}
		}
	}`

	// Проверка результатов с использованием assert.JSONEq
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetInfoHandler_UserIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка запроса с ошибкой получения user_id из контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", "user1") // Устанавливаем только username в контекст

	// Вызов обработчика
	handler.GetInfo(c)

	// Проверка результатов
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неавторизован.")
}

func TestGetInfoHandler_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем ошибку в репозиториях (например, ошибка при получении баланса)
	repository.GetUserBalanceByID = func(userID int) (int, error) {
		return 0, fmt.Errorf("ошибка базы данных")
	}

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1) // Устанавливаем user_id в контекст

	// Вызов обработчика
	handler.GetInfo(c)

	// Проверка результатов
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Внутренняя ошибка сервера.")
}

func TestGetInfoHandler_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем ошибку получения user_id
	repository.GetUserBalanceByID = func(userID int) (int, error) {
		if userID == 1 {
			return 1000, nil
		}
		return 0, fmt.Errorf("баланс не найден")
	}

	// Подготовка запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "invalid") // Неверный тип user_id (не int)

	// Вызов обработчика
	handler.GetInfo(c)

	// Проверка результатов
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неавторизован.")
}
