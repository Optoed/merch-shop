package unit

import (
	"golang.org/x/crypto/bcrypt"
	"merch-shop/pkg/utils"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	// Создаем корректный хэш
	password := "mySecretPassword"
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка при хешировании пароля: %v", err)
	}

	tests := []struct {
		name         string
		password     string
		passwordHash string
		expected     bool
	}{
		{"Correct password", password, passwordHash, true},           // Позитивный тест
		{"Incorrect password", "wrongPassword", passwordHash, false}, // Негативный тест
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.CheckPasswordHash(test.password, test.passwordHash)
			if result != test.expected {
				t.Errorf("CheckPasswordHash(%s, %s) = %v; want %v", test.password, test.passwordHash, result, test.expected)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword"

	// Позитивный тест - хэширование пароля
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка при хешировании пароля: %v", err)
	}

	// Проверяем, что хэшированный пароль не пустой
	if hashedPassword == "" {
		t.Errorf("Ожидали непустую строку")
	}

	// Проверяем, что хэшированный пароль можно расшифровать
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Errorf("Ошибка при сравнении хэша с паролем: %v", err)
	}
}
