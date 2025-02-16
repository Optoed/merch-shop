package service

//
//import (
//	"errors"
//	"merch-shop/internal/infrastructure/database"
//	"merch-shop/internal/repository"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//// MockRepository представляет мок для repository
//type MockRepository struct {
//	mock.Mock
//}
//
//func (m *MockRepository) GetCostByName(itemName string) (int, error) {
//	args := m.Called(itemName)
//	return args.Int(0), args.Error(1)
//}
//
//func (m *MockRepository) GetUserBalanceByID(userID int) (int, error) {
//	args := m.Called(userID)
//	return args.Int(0), args.Error(1)
//}
//
//func (m *MockRepository) DecreaseBalanceByAmountTx(tx interface{}, userID int, amount int) error {
//	args := m.Called(tx, userID, amount)
//	return args.Error(0)
//}
//
//func (m *MockRepository) CheckHaveItemInInventoryTx(tx interface{}, userID int, itemName string) (bool, error) {
//	args := m.Called(tx, userID, itemName)
//	return args.Bool(0), args.Error(1)
//}
//
//func (m *MockRepository) AddItemToInventoryTx(tx interface{}, userID int, itemName string, quantity int) error {
//	args := m.Called(tx, userID, itemName, quantity)
//	return args.Error(0)
//}
//
//func (m *MockRepository) UpdateInventoryItemTx(tx interface{}, userID int, itemName string, quantity int) error {
//	args := m.Called(tx, userID, itemName, quantity)
//	return args.Error(0)
//}
//
//// MockDB представляет мок для базы данных
//type MockDB struct {
//	mock.Mock
//}
//
//func (m *MockDB) Beginx() (interface{}, error) {
//	args := m.Called()
//	return args.Get(0), args.Error(1)
//}
//
//func (m *MockDB) Rollback() error {
//	args := m.Called()
//	return args.Error(0)
//}
//
//func (m *MockDB) Commit() error {
//	args := m.Called()
//	return args.Error(0)
//}
//
//func TestBuyItem_Success(t *testing.T) {
//	// Создаем моки
//	mockRepo := new(MockRepository)
//	mockDB := new(MockDB)
//
//	// Подменяем зависимости
//	repository.Store = mockRepo
//	database.DB = mockDB
//
//	// Настраиваем поведение моков
//	mockRepo.On("GetCostByName", "sword").Return(100, nil)
//	mockRepo.On("GetUserBalanceByID", 1).Return(200, nil)
//	mockDB.On("Beginx").Return(mockDB, nil)
//	mockRepo.On("DecreaseBalanceByAmountTx", mockDB, 1, 100).Return(nil)
//	mockRepo.On("CheckHaveItemInInventoryTx", mockDB, 1, "sword").Return(false, nil)
//	mockRepo.On("AddItemToInventoryTx", mockDB, 1, "sword", 1).Return(nil)
//	mockDB.On("Commit").Return(nil)
//
//	// Выполняем тестируемую функцию
//	err := BuyItem(1, "sword")
//
//	// Проверяем результат
//	assert.NoError(t, err)
//	mockRepo.AssertExpectations(t)
//	mockDB.AssertExpectations(t)
//}
//
//func TestBuyItem_InsufficientBalance(t *testing.T) {
//	// Создаем моки
//	mockRepo := new(MockRepository)
//	mockDB := new(MockDB)
//
//	// Подменяем зависимости
//	repository.Store = mockRepo
//	database.DB = mockDB
//
//	// Настраиваем поведение моков
//	mockRepo.On("GetCostByName", "sword").Return(100, nil)
//	mockRepo.On("GetUserBalanceByID", 1).Return(50, nil)
//
//	// Выполняем тестируемую функцию
//	err := BuyItem(1, "sword")
//
//	// Проверяем результат
//	assert.EqualError(t, err, "Недостаточно монет на счете")
//	mockRepo.AssertExpectations(t)
//	mockDB.AssertExpectations(t)
//}
//
//func TestBuyItem_InvalidItem(t *testing.T) {
//	// Создаем моки
//	mockRepo := new(MockRepository)
//	mockDB := new(MockDB)
//
//	// Подменяем зависимости
//	repository.Store = mockRepo
//	database.DB = mockDB
//
//	// Настраиваем поведение моков
//	mockRepo.On("GetCostByName", "invalid_item").Return(0, errors.New("item not found"))
//
//	// Выполняем тестируемую функцию
//	err := BuyItem(1, "invalid_item")
//
//	// Проверяем результат
//	assert.EqualError(t, err, "Неверный запрос.")
//	mockRepo.AssertExpectations(t)
//	mockDB.AssertExpectations(t)
//}
//
//func TestBuyItem_TransactionError(t *testing.T) {
//	// Создаем моки
//	mockRepo := new(MockRepository)
//	mockDB := new(MockDB)
//
//	// Подменяем зависимости
//	repository.Store = mockRepo
//	database.DB = mockDB
//
//	// Настраиваем поведение моков
//	mockRepo.On("GetCostByName", "sword").Return(100, nil)
//	mockRepo.On("GetUserBalanceByID", 1).Return(200, nil)
//	mockDB.On("Beginx").Return(mockDB, nil)
//	mockRepo.On("DecreaseBalanceByAmountTx", mockDB, 1, 100).Return(errors.New("transaction error"))
//	mockDB.On("Rollback").Return(nil)
//
//	// Выполняем тестируемую функцию
//	err := BuyItem(1, "sword")
//
//	// Проверяем результат
//	assert.EqualError(t, err, "Внутренняя ошибка сервера.")
//	mockRepo.AssertExpectations(t)
//	mockDB.AssertExpectations(t)
//}
