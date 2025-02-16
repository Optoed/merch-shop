package unit

import (
	_ "errors"
	"merch-shop/pkg/config"
	"merch-shop/pkg/jwtAuth"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "merch-shop/pkg/config"
)

// Тестируем функцию GenerateJWT
func TestGenerateJWT(t *testing.T) {
	secretKey := config.Cfg.SecretJWTKey
	userID := 123
	username := "testuser"

	// Генерируем JWT
	token, err := jwtAuth.GenerateJWT(secretKey, userID, username)
	if err != nil {
		t.Fatalf("Ошибка при генерации JWT: %v", err)
	}

	// Парсим JWT и проверяем его
	parsedUserID, parsedUsername, err := jwtAuth.ParseJWT(token)
	if err != nil {
		t.Fatalf("Ошибка при парсинге JWT: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("Expected userID %d, but got %d", userID, parsedUserID)
	}
	if parsedUsername != username {
		t.Errorf("Expected username %s, but got %s", username, parsedUsername)
	}
}

// Тестируем функцию ParseJWT
func TestParseJWT(t *testing.T) {
	// Используем секретный ключ
	secretKey := config.Cfg.SecretJWTKey
	userID := 123
	username := "testuser"

	// Генерируем JWT
	token, err := jwtAuth.GenerateJWT(secretKey, userID, username)
	if err != nil {
		t.Fatalf("Ошибка при генерации JWT: %v", err)
	}

	// Проверка правильного токена
	parsedUserID, parsedUsername, err := jwtAuth.ParseJWT(token)
	if err != nil {
		t.Fatalf("Ошибка при парсинге JWT: %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("Expected userID %d, but got %d", userID, parsedUserID)
	}
	if parsedUsername != username {
		t.Errorf("Expected username %s, but got %s", username, parsedUsername)
	}

	// Тестируем невалидный токен (неправильный секретный ключ)
	invalidToken := "invalid.token.string"
	_, _, err = jwtAuth.ParseJWT(invalidToken)
	if err == nil {
		t.Error("Expected error when parsing invalid token, but got nil")
	}

	// Тестируем токен с истекшим сроком действия
	expiredClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * -24).Unix(), // Истекший токен
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	invalidExpiredToken, _ := expiredToken.SignedString([]byte(secretKey))

	_, _, err = jwtAuth.ParseJWT(invalidExpiredToken)
	if err == nil {
		t.Error("Expected error when parsing expired token, but got nil")
	}
}
