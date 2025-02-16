package unit

import (
	_ "bytes"
	"fmt"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"merch-shop/internal/middleware"
	"merch-shop/pkg/jwtAuth"

	"net/http"
	"net/http/httptest"
	"testing"
)

func mockParseJWT(token string) (int, string, error) {
	if token == "valid_token" {
		return 123, "validUser", nil
	}
	return 0, "", fmt.Errorf("invalid token")
}

func TestJWTMiddleware_NoAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка запроса без заголовка Authorization
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/protected", nil)

	// Применяем middleware
	middleware.JWTMiddleware()(c)

	// Проверяем, что код ответа 401 (Неавторизован)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Неавторизован.")
}

func TestJWTMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Подготовка запроса с неправильным форматом заголовка Authorization
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat token")

	// Применяем middleware
	middleware.JWTMiddleware()(c)

	// Проверяем, что код ответа 401 (Неавторизован)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Неавторизован.")
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем ParseJWT для возврата ошибки
	jwtAuth.ParseJWT = mockParseJWT

	// Подготовка запроса с невалидным токеном
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid_token")

	// Применяем middleware
	middleware.JWTMiddleware()(c)

	// Проверяем, что код ответа 401 (Неавторизован)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Неавторизован.")
}

func TestJWTMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Мокируем ParseJWT для успешного парсинга токена
	jwtAuth.ParseJWT = mockParseJWT

	// Подготовка запроса с валидным токеном
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_token")

	// Применяем middleware
	middleware.JWTMiddleware()(c)

	// Проверяем, что код ответа 200 и установлены данные в контексте
	assert.Equal(t, http.StatusOK, w.Code)
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	assert.Equal(t, 123, userID.(int))
	assert.Equal(t, "validUser", username)
}
