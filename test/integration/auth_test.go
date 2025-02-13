package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"merch-shop/internal/handler"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/models/requestModels"
	"merch-shop/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler(t *testing.T) {

	config.LoadEnv()

	database.InitDB(true) // Передаем true для использования тестовой базы
	database.ClearDB()

	defer func() {
		database.ClearDB()
		database.CloseDB()
	}()

	// Инициализация маршрутов
	router := gin.Default()
	router.POST("/api/auth", handler.AuthHandler)

	t.Run("Success - Auth with correct credentials", func(t *testing.T) {
		// Создаем запрос с правильными данными
		authRequest := requestModels.AuthRequest{
			Username: "testuser",
			Password: "testpassword", // допустим, такой пароль
		}
		requestBody, _ := json.Marshal(authRequest)
		req := bytes.NewReader(requestBody)

		// Выполняем запрос
		w := performRequest(router, "POST", "/api/auth", req)

		// Проверяем успешный ответ
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "token")
	})

	t.Run("Failure - Auth with incorrect credentials", func(t *testing.T) {
		// Создаем запрос с неверными данными
		authRequest := requestModels.AuthRequest{
			Username: "testuser",
			Password: "wrongpassword", // неверный пароль
		}
		requestBody, _ := json.Marshal(authRequest)
		req := bytes.NewReader(requestBody)

		// Выполняем запрос
		w := performRequest(router, "POST", "/api/auth", req)

		// Проверяем, что вернулся статус 401
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "description")
	})
}

// performRequest - вспомогательная функция для выполнения HTTP запроса
func performRequest(router *gin.Engine, method, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
