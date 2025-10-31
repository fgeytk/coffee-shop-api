package handlers

import (
	"coffee-shop-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var orders []models.Order
var orderCounter int = 1

func calculatePrice(basePrice float64, size string, extras []string) float64 {
	//prix de base
	total := basePrice

	//taille
	switch size {
	case "small":
		total *= 0.8
	case "medium":
		total *= 1.0 // Prix normal
	case "large":
		total *= 1.3
	}

	//prix des extras
	total += float64(len(extras)) * 0.50

	//
	return float64(int(total*100)) / 100
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	//gerer l'erreur
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Données de commande invalides"})
		return
	}
	//verifie si la boisson existe
	var foundDrink *models.Drink
	for _, drink := range models.Drinks {
		if drink.ID == newOrder.DrinkID {
			foundDrink = &drink
			break
		}
	}
	if foundDrink == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Boisson introuvable"})
		return
	}
	// identifiant unique
	newOrder.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++

	newOrder.DrinkName = foundDrink.Name
	newOrder.Status = models.StatusPending
	newOrder.OrderedAt = time.Now()

	// prix (a implémenter calculatePrice)
	newOrder.TotalPrice = calculatePrice(foundDrink.BasePrice, newOrder.Size, newOrder.Extras)

	orders = append(orders, newOrder)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
	fmt.Printf("Nouvelle commande: %+v\n", newOrder)

}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
	fmt.Println("Liste des commandes :")
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]
	for _, order := range orders {
		if order.ID == orderId {
			json.NewEncoder(w).Encode(order)
			fmt.Printf("Détails de la commande %s : %+v\n", orderId, order)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]
	var statusUpdate struct {
		Status models.OrderStatus `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Données de mise à jour invalides"})
		return
	}
	for i, order := range orders {
		if order.ID == orderId {
			orders[i].Status = statusUpdate.Status
			json.NewEncoder(w).Encode(orders[i])
			fmt.Printf("Mise à jour de la commande %s : %+v\n", orderId, orders[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]

	for i, order := range orders {
		if order.ID == orderId {
			if order.Status == models.StatusPickedUp {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Impossible d'annuler une commande déjà récupérée"})
				return
			}
			orders = append(orders[:i], orders[i+1:]...)
			//204
			w.WriteHeader(http.StatusNoContent)
			fmt.Printf("Commande %s supprimée\n", orderId)
			return
		}
	}

	//404
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
}
