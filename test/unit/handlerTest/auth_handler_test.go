package unit

//
//import (
//	"bytes"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"merch-shop/internal/handler"
//	"merch-shop/internal/infrastructure/database"
//	_ "merch-shop/internal/models/requestModels"
//	"net/http"
//	"net/http/httptest"
//	"regexp"
//	"testing"
//
//	// "merch-shop/internal/repositories"
//	_ "merch-shop/internal/service"
//)
//
//func TestAuthHandler(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	// Инициализация мокированной базы данных
//	sqlxDB, mock := setupMockDB()
//	database.DB = sqlxDB
//
//	// Мокируем поведение при проверке пользователя: если пользователя нет в базе, то он должен быть добавлен
//	mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE username=?")).
//		WithArgs("validUser").
//		WillReturnRows(sqlmock.NewRows([]string{"username"})) // Пользователь не найден, добавляем нового
//
//	// Мокируем добавление нового пользователя в базу данных
//	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users")).
//		WithArgs("validUser", sqlmock.AnyArg(), sqlmock.AnyArg()).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// Мокируем генерацию токена
//	token := "mockedJWTToken"
//
//	// Создаем новый HTTP-запрос для теста
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//
//	// Формируем тело запроса
//	requestBody := bytes.NewBufferString(`{"username": "validUser", "password": "correctPassword"}`)
//	c.Request, _ = http.NewRequest("POST", "/auth", requestBody)
//	c.Request.Header.Set("Content-Type", "application/json")
//
//	// Вызов обработчика
//	handler.AuthHandler(c)
//
//	// Проверка на успешную аутентификацию
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Contains(t, w.Body.String(), "Успешная аутентификация.")
//	assert.Contains(t, w.Body.String(), token) // Проверка на сгенерированный токен
//
//	// Проверка, что все ожидания для mock были выполнены
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestAuthHandler_Unauthorized(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	// Инициализация мокированной базы данных
//	sqlxDB, mock := setupMockDB()
//	database.DB = sqlxDB
//
//	// Мокируем поведение при проверке пользователя: если пользователя нет в базе, то он должен быть добавлен
//	mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE username=?")).
//		WithArgs("invalidUser").
//		WillReturnRows(sqlmock.NewRows([]string{"username"})) // Пользователь не найден в базе
//
//	// Создаем новый HTTP-запрос для теста
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//
//	// Формируем тело запроса
//	requestBody := bytes.NewBufferString(`{"username": "invalidUser", "password": "wrongPassword"}`)
//	c.Request, _ = http.NewRequest("POST", "/auth", requestBody)
//	c.Request.Header.Set("Content-Type", "application/json")
//
//	// Вызов обработчика
//	handler.AuthHandler(c)
//
//	// Проверка на ошибку 401 (Не авторизован)
//	assert.Equal(t, http.StatusUnauthorized, w.Code)
//	assert.Contains(t, w.Body.String(), "Неавторизован")
//
//	// Проверка, что все ожидания для mock были выполнены
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
