package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/models"
)

func GetUserInventory(userID int) ([]models.InventoryItemResponse, error) {
	query := `SELECT item_name, count FROM inventory WHERE user_id=$1`
	var inventory []models.InventoryItemResponse
	err := database.DB.Select(&inventory, query, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ошибка при получении инвентаря пользователя: %v", err)
		return nil, err
	}
	return inventory, nil
}

func CheckHaveItemInInventoryTx(tx *sqlx.Tx, userID int, itemName string) (bool, error) {
	query := `SELECT 1 FROM inventory WHERE user_id=$1 AND item_name=$2 LIMIT 1`
	var exists bool
	err := tx.Get(&exists, query, userID, itemName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists, nil
}

func CheckHaveItemInInventory(userID int, itemName string) (bool, error) {
	query := `SELECT 1 FROM inventory WHERE user_id=$1 AND item_name=$2 LIMIT 1`
	var exists bool
	err := database.DB.Get(&exists, query, userID, itemName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists, nil
}

func AddItemToInventoryTx(tx *sqlx.Tx, userID int, itemName string, count int) error {
	query := `INSERT INTO inventory (user_id, item_name, count) VALUES ($1, $2, $3)`
	_, err := tx.Exec(query, userID, itemName, count)
	if err != nil {
		log.Printf("Ошибка при добавлении товара в инвентарь: %v", err)
		return err
	}
	return nil
}

func UpdateInventoryItemTx(tx *sqlx.Tx, userID int, itemName string, count int) error {
	query := `UPDATE inventory SET count = count + $1 WHERE user_id = $2 AND item_name = $3`
	_, err := tx.Exec(query, count, userID, itemName)
	if err != nil {
		log.Printf("Ошибка при обновлении инвентаря: %v", err)
		return err
	}
	return nil
}
