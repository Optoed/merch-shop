package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"merch-shop/internal/models/requestModels"
	"net/http"
	"net/http/httptest"
	"testing"
)

// PerformRequest - вспомогательная функция для выполнения HTTP запроса
func PerformRequest(router *gin.Engine, method, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// AuthenticateUser - вспомогательная функция для аутентификации пользователя и получения токена
func AuthenticateUser(t *testing.T, router *gin.Engine, username, password string) string {
	authRequest := requestModels.AuthRequest{
		Username: username,
		Password: password,
	}
	requestBody, _ := json.Marshal(authRequest)
	req := bytes.NewReader(requestBody)

	w := PerformRequest(router, "POST", "/api/auth", req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"].(string)
}

// PerformAuthenticatedRequest - вспомогательная функция для выполнения аутентифицированного HTTP запроса
func PerformAuthenticatedRequest(router *gin.Engine, method, path string, body *bytes.Reader, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
