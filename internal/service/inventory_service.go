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
		return err
	}

	userBalance, err := repository.GetUserBalanceByID(userID)
	if err != nil {
		return err
	}

	if userBalance < cost {
		return errors.New("Недостаточно монет на счете")
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v", err)
		return err
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
		return err
	}

	//TODO - допиши
	//
	_, err = repository.GetUserInventory(userID)
	if err != nil {
		err = repository.AddItemToInventory(userID, itemName, 1)
		if err != nil {
			return err
		}
	} else {
		err = repository.UpdateInventoryItem(userID, itemName, 1)
		if err != nil {
			return err
		}
	}
	//

	if err = tx.Commit(); err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return err
	}

	return nil
}
