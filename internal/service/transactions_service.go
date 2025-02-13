package service

import (
	"errors"
	"log"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/repository"
)

func SendCoin(senderID int, receiverName string, amount int) error {
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

	senderBalance, err := repository.GetBalanceByID(senderID)
	//log.Printf("senderID = %d, senderBalance = %d, error = %v",
	//	senderID, senderBalance, err)
	if err != nil {
		return err
	}

	if senderBalance < amount {
		return errors.New("Недостаточно монет на счете")
	}

	receiver, err := repository.GetUserByUsername(receiverName)
	if err != nil {
		return err
	}

	//TODO нужно ли?
	//if senderID == receiver.ID {
	//	return errors.New("Нельзя поделиться монетами с самим же собой")
	//}

	err = repository.UpdateBalanceTx(tx, senderID, -amount)
	if err != nil {
		return err
	}

	err = repository.UpdateBalanceTx(tx, receiver.ID, +amount)
	if err != nil {
		return err
	}

	err = repository.CreateTransactionTx(tx, senderID, receiver.ID, receiverName, amount)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return err
	}

	return nil
}
