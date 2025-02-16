package integration

import (
	_ "bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
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

func TestGetUserInfo(t *testing.T) {
	// Загружаем переменные окружения
	config.LoadEnv()

	// Инициализируем тестовую базу данных
	database.InitDB(true)
	database.ClearDB()

	defer func() {
		database.ClearDB()
	}()

	// Инициализация маршрутов
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Подготовка для создания пользователя и авторизации
	router.POST("/api/auth", handler.AuthHandler)
	protectedAPI := router.Group("/api")
	protectedAPI.Use(middleware.JWTMiddleware())
	{
		protectedAPI.GET("/info", handler.GetInfo)
	}

	testUsername := "testuser"
	testPassword := "testpassword"

	testUsername2 := "anotherUser"
	testPassword2 := "anotherPassword"

	// Создаем тестового пользователя
	userID, err := repository.CreateUser(testUsername, testPassword, 1000)
	assert.NoError(t, err)

	userID2, _ := repository.CreateUser(testUsername2, testPassword2, 1000)

	// Авторизуем пользователя
	token := test.AuthenticateUser(t, router, testUsername, testPassword)

	// Тест на получение информации о пользователе
	t.Run("Success - Get user info", func(t *testing.T) {

		item := "cup"

		// Добавим в инвентарь 1 предмет "cup"
		query := `INSERT INTO inventory (user_id, item_name, count) VALUES ($1, $2, $3)`
		_, err := database.DB.Exec(query, userID, item, 1)
		assert.NoError(t, err)

		// Добавим 1 транзакцию
		query = `INSERT INTO transactions (sender_id, sender_name, receiver_id, receiver_name, amount) 
			  VALUES ($1, $2, $3, $4, $5)`
		_, err = database.DB.Exec(query, userID, testUsername, userID2, testUsername2, 50)
		assert.NoError(t, err)

		// Выполняем запрос на получение информации
		getUserInfoReq, _ := http.NewRequest("GET", "/api/info", nil)
		getUserInfoReq.Header.Set("Authorization", "Bearer "+token)
		getUserInfoW := httptest.NewRecorder()

		// Выполняем запрос через сервер
		router.ServeHTTP(getUserInfoW, getUserInfoReq)

		// Проверяем успешный статус код
		assert.Equal(t, http.StatusOK, getUserInfoW.Code)

		// Проверяем, что информация о пользователе корректна
		var response map[string]interface{}
		json.Unmarshal(getUserInfoW.Body.Bytes(), &response)

		log.Println("response = ", response)

		// Проверяем что в ответе правильное описание
		assert.Equal(t, "Успешный ответ.", response["description"])

		// Проверяем, что в ответе есть инвентарь и транзакции
		schema := response["schema"].(map[string]interface{})
		inventory := schema["inventory"].([]interface{})
		coinHistory := schema["coinHistory"].(map[string]interface{})
		coins := schema["coins"]

		assert.Equal(t, float64(1000), coins)

		// Проверяем, что инвентарь не пустой и содержит предмет "cup"
		assert.NotEmpty(t, inventory)
		assert.Equal(t, "cup", inventory[0].(map[string]interface{})["type"].(string))
		assert.Equal(t, float64(1), inventory[0].(map[string]interface{})["quantity"])

		//// Проверяем, что транзакции не пустые
		received := coinHistory["received"]
		assert.Equal(t, nil, received)

		sent := coinHistory["sent"].([]interface{})
		assert.NotEmpty(t, sent)

		// Проверяем, что транзакция имеет правильные значения
		assert.Equal(t, 50.0, sent[0].(map[string]interface{})["amount"])
		assert.Equal(t, testUsername2, sent[0].(map[string]interface{})["toUser"])
	})
}
