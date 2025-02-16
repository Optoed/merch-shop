package unit

import (
	"database/sql"
	"errors"
	_ "errors"
	"merch-shop/internal/models"
	"merch-shop/pkg/config"
	"merch-shop/pkg/jwtAuth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merch-shop/pkg/utils"
)

// Мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) CreateUser(username, password string, balance int) (int, error) {
	args := m.Called(username, password, balance)
	return args.Int(0), args.Error(1)
}

// Мок для утилиты проверки пароля
func mockCheckPasswordHash(password, hash string) bool {
	return password == hash
}

func Authenticate(repo *MockRepository, username, password string) (string, error) {
	user, err := repo.GetUserByUsername(username)

	if errors.Is(err, sql.ErrNoRows) {
		return createUserAndGenerateJWT(repo, username, password)
	} else if err != nil {
		return "", errors.New("Неавторизован.")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("Неавторизован.")
	}

	return generateJWT(user.ID, username)
}

func createUserAndGenerateJWT(repo *MockRepository, username, password string) (string, error) {
	userID, err := repo.CreateUser(username, password, 1000)
	if err != nil {
		return "", errors.New("внутренняя ошибка сервера")
	}
	return generateJWT(userID, username)
}

func generateJWT(userID int, username string) (string, error) {
	token, err := jwtAuth.GenerateJWT(config.Cfg.SecretJWTKey, userID, username)
	if err != nil {
		return "", errors.New("внутренняя ошибка сервера")
	}
	return token, nil
}

func TestAuthenticate(t *testing.T) {
	// Мокаем репозиторий
	mockRepo := new(MockRepository)
	utils.CheckPasswordHash = mockCheckPasswordHash // Мокаем CheckPasswordHash

	t.Run("User does not exist, creates user and returns JWT", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "newuser").Return(nil, sql.ErrNoRows)
		mockRepo.On("CreateUser", "newuser", "password123", 1000).Return(1, nil)

		jwtToken, err := Authenticate(mockRepo, "newuser", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, jwtToken)
		mockRepo.AssertExpectations(t)
	})

	//t.Run("User exists and password is correct, returns JWT", func(t *testing.T) {
	//	mockRepo.On("GetUserByUsername", "existinguser").Return(&models.User{
	//		ID:           1,
	//		Username:     "existinguser",
	//		PasswordHash: "password123", // Подразумеваем, что пароль хэширован правильно
	//	}, nil)
	//	jwtToken, err := service.Authenticate("existinguser", "password123")
	//
	//	assert.NoError(t, err)
	//	assert.NotEmpty(t, jwtToken)
	//	mockRepo.AssertExpectations(t)
	//})
	//
	//t.Run("Password is incorrect, returns error", func(t *testing.T) {
	//	mockRepo.On("GetUserByUsername", "existinguser").Return(&repository.User{
	//		ID:           1,
	//		Username:     "existinguser",
	//		PasswordHash: "incorrectpassword",
	//	}, nil)
	//	jwtToken, err := Authenticate("existinguser", "wrongpassword")
	//
	//	assert.Error(t, err)
	//	assert.Empty(t, jwtToken)
	//	mockRepo.AssertExpectations(t)
	//})
}
