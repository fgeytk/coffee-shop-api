package models

type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"` // coffee, tea, cold
	BasePrice float64 `json:"base_price"`
}

// Base de données en mémoire
var Drinks = []Drink{ // Exported variable for access in main.go
	{ID: "1", Name: "Espresso", Category: "coffee", BasePrice: 2.0},
	{ID: "2", Name: "Cappuccino", Category: "coffee", BasePrice: 3.0},
	{ID: "3", Name: "Latte", Category: "coffee", BasePrice: 3.5},
	{ID: "4", Name: "Black Tea", Category: "tea", BasePrice: 2.5},
	{ID: "5", Name: "Green Tea", Category: "tea", BasePrice: 2.5},
	{ID: "6", Name: "Iced Coffee", Category: "cold", BasePrice: 3.0},
	{ID: "7", Name: "Iced Tea", Category: "cold", BasePrice: 2.5},
}
