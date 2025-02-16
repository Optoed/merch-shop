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

	// Создаем тестового пользователя
	userID, err := repository.CreateUser(testUsername, testPassword, 1000)
	assert.NoError(t, err)

	anotherUserID, _ := repository.CreateUser("anotherUser", "anotherPassword", 1000)

	// Авторизуем пользователя
	token := test.AuthenticateUser(t, router, testUsername, testPassword)

	// Тест на получение информации о пользователе
	t.Run("Success - Get user info", func(t *testing.T) {

		// Добавим в инвентарь 1 предмет "cup"
		query := `INSERT INTO inventory (user_id, item_name, count) VALUES ($1, $2, $3)`
		_, err := database.DB.Exec(query, userID, "cup", 1)
		assert.NoError(t, err)

		// Добавим 1 транзакцию
		query = `INSERT INTO transactions (sender_id, receiver_id, receiver_name, amount) 
			  VALUES ($1, $2, $3, $4)`
		_, err = database.DB.Exec(query, userID, anotherUserID, "anotherUser", 50)
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
		assert.Equal(t, 1000, int(response["schema"].(map[string]interface{})["coins"].(float64)))

		// Проверяем, что в ответе есть инвентарь и транзакции
		schema := response["schema"].(map[string]interface{})
		inventory := schema["inventory"].([]interface{})
		coinHistory := schema["coinHistory"].(map[string]interface{})

		// Проверяем, что инвентарь не пустой и содержит предмет "cup"
		assert.NotEmpty(t, inventory)
		assert.Equal(t, "cup", inventory[0].(map[string]interface{})["item_name"])
		assert.Equal(t, float64(1), inventory[0].(map[string]interface{})["count"])

		// Проверяем, что транзакции не пустые
		received := coinHistory["received"].([]interface{})
		sent := coinHistory["sent"].([]interface{})
		assert.NotEmpty(t, received)
		assert.NotEmpty(t, sent)

		// Проверяем, что транзакция имеет правильные значения
		assert.Equal(t, 50.0, received[0].(map[string]interface{})["amount"])
		assert.Equal(t, "anotherUser", received[0].(map[string]interface{})["sender_name"])
		assert.Equal(t, 50.0, sent[0].(map[string]interface{})["amount"])
		assert.Equal(t, "anotherUser", sent[0].(map[string]interface{})["receiver_name"])
	})

	//t.Run("Failure - Get user info without token", func(t *testing.T) {
	//	// Выполняем запрос без авторизации
	//	getUserInfoReq, _ := http.NewRequest("GET", "/api/info", nil)
	//	getUserInfoW := httptest.NewRecorder()
	//
	//	// Выполняем запрос через сервер
	//	router.ServeHTTP(getUserInfoW, getUserInfoReq)
	//
	//	// Проверяем, что вернулся статус 401 (не авторизован)
	//	assert.Equal(t, http.StatusUnauthorized, getUserInfoW.Code)
	//
	//	// Проверяем, что в ответе есть описание ошибки
	//	var response map[string]interface{}
	//	json.Unmarshal(getUserInfoW.Body.Bytes(), &response)
	//	assert.Contains(t, response, "description")
	//})
}
