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
	"net/http"
	"testing"
)

func TestBuyItemHandler(t *testing.T) {

	config.LoadEnv()

	database.InitDB(true) // Передаем true для использования тестовой базы
	database.ClearDB()

	defer func() {
		database.ClearDB()
	}()

	gin.SetMode(gin.TestMode)

	// Инициализация маршрутов
	router := gin.Default()
	router.POST("/api/auth", handler.AuthHandler)
	protectedAPI := router.Group("/api")
	protectedAPI.Use(middleware.JWTMiddleware())
	{
		protectedAPI.POST("/buy/:item", handler.BuyItem)
	}

	testUsername := "testuser"
	testPassword := "testpassword"

	// Создаем тестового пользователя
	userID, err := repository.CreateUser(testUsername, testPassword, 1000)
	assert.NoError(t, err)

	token := AuthenticateUser(t, router, testUsername, testPassword)
	//log.Printf("token = %s\n", token)

	t.Run("Success - Buy item with sufficient balance", func(t *testing.T) {

		//log.Println("db = ", database.DB)

		// Выполняем запрос на покупку товара
		itemName := "cup"
		itemCost, err := repository.Store.GetCostByName(itemName)
		assert.NoError(t, err)

		requestBody, _ := json.Marshal(itemName)
		req := bytes.NewReader(requestBody)
		w := PerformAuthenticatedRequest(router, "POST", "/api/buy/"+itemName, req, token)

		// Проверяем успешный ответ
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Успешный ответ.", response["description"])

		// Проверяем, что баланс пользователя уменьшился
		userBalance, err := repository.GetUserBalanceByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, 1000-itemCost, userBalance)

		// Проверяем, что товар добавлен в инвентарь
		exists, err := repository.CheckHaveItemInInventory(userID, itemName)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Failure - Buy item with insufficient balance", func(t *testing.T) {
		// Устанавливаем баланс пользователя меньше стоимости товара

		itemName := "pink-hoody"
		itemCost, err := repository.Store.GetCostByName(itemName)

		err = repository.SetUserBalance(userID, itemCost-1)
		assert.NoError(t, err)

		// Выполняем запрос на покупку товара
		requestBody, _ := json.Marshal(itemName)
		req := bytes.NewReader(requestBody)
		w := PerformAuthenticatedRequest(router, "POST", "/api/buy/"+itemName, req, token)

		// Проверяем, что вернулся статус 400
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Недостаточно монет на счете", response["description"])
	})
}
