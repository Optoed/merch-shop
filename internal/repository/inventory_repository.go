package repository

import (
	"log"
	"merch-shop/internal/infrastructure/database"
	"merch-shop/internal/models"
)

func GetUserInventory(userID int) ([]models.InventoryItem, error) {
	query := `SELECT * FROM inventory WHERE user_id=$1`
	var inventory []models.InventoryItem
	err := database.DB.Select(&inventory, query, userID)
	if err != nil {
		log.Printf("Ошибка при получении инвентаря пользователя: %v", err)
		return nil, err
	}
	return inventory, nil
}

func AddItemToInventory(userID int, itemName string, count int) error {
	query := `INSERT INTO inventory (user_id, item_name, count) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, userID, itemName, count)
	if err != nil {
		log.Printf("Ошибка при добавлении товара в инвентарь: %v", err)
		return err
	}
	return nil
}

func UpdateInventoryItem(userID int, itemName string, count int) error {
	query := `UPDATE inventory SET count = count + $1 WHERE user_id = $2 AND item_name = $3`
	_, err := database.DB.Exec(query, count, userID, itemName)
	if err != nil {
		log.Printf("Ошибка при обновлении инвентаря: %v", err)
		return err
	}
	return nil
