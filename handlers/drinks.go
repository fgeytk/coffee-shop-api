package handlers

import (
	"coffee-shop-api/database"
	"coffee-shop-api/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT id, name, category, base_price FROM drinks")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la récupération du menu"})
		return
	}
	defer rows.Close()

	var drinks []models.Drink
	for rows.Next() {
		var drink models.Drink
		err := rows.Scan(&drink.ID, &drink.Name, &drink.Category, &drink.BasePrice)
		if err != nil {
			continue
		}
		drinks = append(drinks, drink)
	}

	json.NewEncoder(w).Encode(drinks)
	fmt.Println("Menu récupéré depuis MySQL")
}

func GetDrink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	drinkID := vars["id"]

	var drink models.Drink
	err := database.DB.QueryRow("SELECT id, name, category, base_price FROM drinks WHERE id = ?", drinkID).
		Scan(&drink.ID, &drink.Name, &drink.Category, &drink.BasePrice)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Boisson introuvable"})
		return
	}

	json.NewEncoder(w).Encode(drink)
}
