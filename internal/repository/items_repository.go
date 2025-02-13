package repository

import "errors"

type Items struct {
	dict map[string]int
}

var Store Items = Items{dict: map[string]int{
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
}}

func (items *Items) GetCostByName(name string) (int, error) {
	if cost, exists := items.dict[name]; !exists {
		return 0, errors.New("Товара с таким названием не существует")
	} else {
		return cost, nil
	}
}
