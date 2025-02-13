package repository

import "errors"

type ItemsMap map[string]int

var Store = ItemsMap{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

func (items ItemsMap) GetCostByName(name string) (int, error) {
	if cost, exists := items[name]; !exists {
		return 0, errors.New("Товара с таким названием не существует")
	} else {
		return cost, nil
	}
}
