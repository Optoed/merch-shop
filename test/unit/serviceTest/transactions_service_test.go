package unit

//import (
//	"merch-shop/internal/infrastructure/database"
//	"merch-shop/internal/service"
//	"regexp"
//	"testing"
//
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/jmoiron/sqlx"
//	"github.com/stretchr/testify/assert"
//)
//
//func setupMockDB() (*sqlx.DB, sqlmock.Sqlmock) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		panic(err)
//	}
//	sqlxDB := sqlx.NewDb(db, "sqlmock")
//	return sqlxDB, mock
//}
//
//func TestSendCoinService(t *testing.T) {
//	sqlxDB, mock := setupMockDB()
//	database.DB = sqlxDB
//
//	// Мокируем получение баланса отправителя
//	mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM users WHERE id=$1")).
//		WithArgs(1).
//		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(1000)) // У пользователя достаточно монет
//
//	// Мокируем получение ID получателя по имени
//	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE username=$1")).
//		WithArgs("user2").
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2)) // Получаем ID пользователя
//
//	// Начинаем транзакцию
//	mock.ExpectBegin()
//
//	// Мокируем уменьшение монет на счету отправителя
//	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET balance = balance - $1 WHERE id = $2")).
//		WithArgs(100, 1).                         // Отправляем 100 монет от пользователя с ID 1
//		WillReturnResult(sqlmock.NewResult(0, 1)) // Ожидаем успешное обновление
//
//	// Мокируем увеличение монет на счету получателя
//	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET balance = balance + $1 WHERE id = $2")).
//		WithArgs(100, 2).                         // Получаем 100 монет на счет пользователя с ID 2
//		WillReturnResult(sqlmock.NewResult(0, 1)) // Ожидаем успешное обновление
//
//	// Мокируем создание транзакции
//	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions")).
//		WithArgs(1, "user1", 2, "user2", 100).    // Параметры транзакции
//		WillReturnResult(sqlmock.NewResult(1, 1)) // Ожидаем успешную вставку
//
//	// Мокируем commit транзакции
//	mock.ExpectCommit()
//
//	// Вызов функции SendCoin
//	err := service.SendCoin(1, "user1", "user2", 100)
//
//	// Проверка результатов
//	assert.NoError(t, err)                        // Ошибка не должна возникнуть
//	assert.NoError(t, mock.ExpectationsWereMet()) // Проверка выполнения всех ожидаемых запросов
//}
//
////func TestSendCoinService_InsufficientBalance(t *testing.T) {
////	sqlxDB, mock := setupMockDB()
////	repositories.DB = sqlxDB
////
////	// Мокируем получение баланса отправителя
////	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1")).
////		WithArgs(1).
////		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(50)) // У пользователя недостаточно монет
////
////	// Вызов функции SendCoin
////	err := SendCoin(1, "user1", "user2", 100)
////
////	// Проверка ошибок
////	assert.EqualError(t, err, "Неверный запрос.") // Ошибка недостаточно монет
////}
////
////func TestSendCoinService_SenderEqualsReceiver(t *testing.T) {
////	sqlxDB, mock := setupMockDB()
////	repositories.DB = sqlxDB
////
////	// Мокируем получение баланса отправителя
////	mock.ExpectQuery(regexp.QuoteMeta("SELECT coins FROM users WHERE id=$1")).
////		WithArgs(1).
////		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000)) // У пользователя достаточно монет
////
////	// Мокируем получение ID получателя по имени (тот же пользователь)
////	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE name=$1")).
////		WithArgs("user1").
////		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Получаем ID того же пользователя
////
////	// Вызов функции SendCoin
////	err := SendCoin(1, "user1", "user1", 100)
////
////	// Проверка ошибок
////	assert.EqualError(t, err, "Неверный запрос.") // Ошибка, нельзя отправлять себе монеты
////}
