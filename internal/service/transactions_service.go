package service

import (
	"errors"
	"log"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/repository"
)

var SendCoin = func(senderID int, senderName string, receiverName string, amount int) error {
	if amount <= 0 {
		return errors.New("Неверный запрос.") // Количество отправляемых монет должно быть положительным!
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v", err)
		return errors.New("Внутренняя ошибка сервера.")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
		}
	}()

	senderBalance, err := repository.GetUserBalanceByID(senderID)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	if senderBalance < amount {
		return errors.New("Неверный запрос.") //Недостаточно монет на счете
	}

	receiver, err := repository.GetUserByUsername(receiverName)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	if senderID == receiver.ID {
		return errors.New("Неверный запрос.") //Нельзя поделиться монетами с самим же собой
	}

	err = repository.DecreaseBalanceByAmountTx(tx, senderID, amount)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	err = repository.IncreaseBalanceByAmountTx(tx, receiver.ID, amount)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	err = repository.CreateTransactionTx(tx, senderID, senderName, receiver.ID, receiverName, amount)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return errors.New("Внутренняя ошибка сервера.")
	}

	return nil
}
