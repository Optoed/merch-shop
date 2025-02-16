package service

import (
	"errors"
	"log"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/repository"
)

func BuyItem(userID int, itemName string) error {
	cost, err := repository.Store.GetCostByName(itemName)
	if err != nil {
		return errors.New("Неверный запрос.")
	}

	userBalance, err := repository.GetUserBalanceByID(userID)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	if userBalance < cost {
		return errors.New("Недостаточно монет на счете")
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

	err = repository.DecreaseBalanceByAmountTx(tx, userID, cost)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	exists, err := repository.CheckHaveItemInInventoryTx(tx, userID, itemName)
	if err != nil {
		return errors.New("Внутренняя ошибка сервера.")
	}

	if !exists {
		err = repository.AddItemToInventoryTx(tx, userID, itemName, 1)
		if err != nil {
			return errors.New("Внутренняя ошибка сервера.")
		}
	} else {
		err = repository.UpdateInventoryItemTx(tx, userID, itemName, 1)
		if err != nil {
			return errors.New("Внутренняя ошибка сервера.")
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return errors.New("Внутренняя ошибка сервера.")
	}

	return nil
}
