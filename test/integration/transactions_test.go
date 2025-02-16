package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"merch-shop/internal/handler"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/middleware"
	"merch-shop/internal/repository"
	"merch-shop/pkg/config"
	"merch-shop/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendCoinHandler(t *testing.T) {

	config.LoadEnv()

	// Инициализируем тестовую базу данных
	database.InitDB(true) // Передаем true для использования тестовой базы
	database.ClearDB()

	defer func() {
		database.ClearDB()
	}()

	gin.SetMode(gin.TestMode)

	// Инициализация маршрутов
	router := gin.Default()
	router.POST("/api/auth", handler.AuthHandler)

	// Защищенный API с JWT middleware
	protectedAPI := router.Group("/api")
	protectedAPI.Use(middleware.JWTMiddleware())
	{
		protectedAPI.POST("/sendCoin", handler.SendCoinHandler)
	}

	// Создаем пользователей для теста
	senderUsername := "sender"
	receiverUsername := "receiver"
	senderPassword := "password123"
	receiverPassword := "password123"

	// Регистрация пользователей
	senderUserID, err := repository.CreateUser(senderUsername, senderPassword, 1000)
	assert.NoError(t, err)
	receiverUserID, err := repository.CreateUser(receiverUsername, receiverPassword, 1000)
	assert.NoError(t, err)

	// Авторизация пользователей
	senderToken := test.AuthenticateUser(t, router, senderUsername, senderPassword)

	t.Run("Success - Transfer coins from sender to receiver", func(t *testing.T) {
		// Подготовка запроса на перевод монет
		sendCoinRequest := map[string]interface{}{
			"receiver_name": "receiver",
			"amount":        50,
		}

		sendCoinRequestBody, _ := json.Marshal(sendCoinRequest)
		sendCoinReq, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendCoinRequestBody))
		sendCoinReq.Header.Set("Authorization", "Bearer "+senderToken)
		sendCoinW := httptest.NewRecorder()

		// Выполняем запрос на перевод монет
		router.ServeHTTP(sendCoinW, sendCoinReq)

		// Проверяем статус код
		assert.Equal(t, http.StatusOK, sendCoinW.Code)

		// Проверяем, что монеты были переведены
		sender, err := repository.GetUserByID(senderUserID)
		assert.NoError(t, err)
		receiver, err := repository.GetUserByID(receiverUserID)
		assert.NoError(t, err)

		// Проверяем, что монеты у отправителя уменьшились
		assert.Equal(t, 950, sender.Balance) // 1000 - 50 = 950
		// Проверяем, что монеты у получателя увеличились
		assert.Equal(t, 1050, receiver.Balance) // 1000 + 50 = 1050

		// Проверяем, что транзакция была записана
		transactions, err := repository.GetTransactionsFromUser(senderUserID)
		assert.NoError(t, err)
		assert.Equal(t, senderUserID, transactions[0].SenderName)
		assert.Equal(t, 50, transactions[0].Amount)
	})

	t.Run("Failure - Insufficient balance", func(t *testing.T) {
		// Подготовим запрос, чтобы отправитель не мог перевести больше монет, чем у него есть
		sendCoinRequest := map[string]interface{}{
			"receiver_name": "receiver",
			"amount":        2000, // Попытка перевести 2000 монет, а у отправителя всего 1000
		}
		sendCoinRequestBody, _ := json.Marshal(sendCoinRequest)
		sendCoinReq, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendCoinRequestBody))
		sendCoinReq.Header.Set("Authorization", "Bearer "+senderToken)
		sendCoinW := httptest.NewRecorder()

		// Выполняем запрос на перевод монет
		router.ServeHTTP(sendCoinW, sendCoinReq)

		// Проверяем статус код
		assert.Equal(t, http.StatusBadRequest, sendCoinW.Code)

		// Проверяем, что у отправителя и получателя монеты не изменились
		sender, err := repository.GetUserByID(senderUserID)
		assert.NoError(t, err)
		receiver, err := repository.GetUserByID(receiverUserID)
		assert.NoError(t, err)

		// Убедимся, что баланс отправителя и получателя не изменился
		assert.Equal(t, 1000, sender.Balance)
		assert.Equal(t, 1000, receiver.Balance)
	})
}
