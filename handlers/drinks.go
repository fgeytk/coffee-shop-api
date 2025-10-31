package handlers

import (
	"coffee-shop-api/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Drinks)
	fmt.Println("Menu :")
}

func GetDrink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	drinkID := vars["id"]
	for _, drink := range models.Drinks {
		if drink.ID == drinkID {
			json.NewEncoder(w).Encode(drink)
			return
		}
	}
}
