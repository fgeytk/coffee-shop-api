package main

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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Drinks) // Use the Drinks slice from models
	fmt.Println("Menu :")
}

func getDrink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	drinkID := vars["id"]
	for _, drink := range models.Drinks { // Use the Drinks slice from models
		if drink.ID == drinkID {
			json.NewEncoder(w).Encode(drink)
			return
		}
	}
}

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

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newOrder models.Order // Use the Order struct from models
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	//gerer l'erreur
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Données de commande invalides"})
		return
	}
	//verifie si la boisson existe
	var foundDrink *models.Drink
	for _, drink := range models.Drinks { // Use the Drinks slice from models
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
	newOrder.Status = models.StatusPending // Use the StatusPending from models
	newOrder.OrderedAt = time.Now()

	// prix (a implémenter calculatePrice)
	newOrder.TotalPrice = calculatePrice(foundDrink.BasePrice, newOrder.Size, newOrder.Extras)

	orders = append(orders, newOrder)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
	fmt.Printf("Nouvelle commande: %+v\n", newOrder)

}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
	fmt.Println("Liste des commandes :")
}

func getOrder(w http.ResponseWriter, r *http.Request) {
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

func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]
	var statusUpdate struct {
		Status models.OrderStatus `json:"status"` // Use the OrderStatus from models
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

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]

	for i, order := range orders {
		if order.ID == orderId {
			if order.Status == models.StatusPickedUp { // Use the StatusPickedUp from models
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

func main() {

	//Routeur mux
	r := mux.NewRouter()

	//Def les routes
	r.HandleFunc("/menu", getMenu).Methods("GET")
	r.HandleFunc("/drinks/{id}", getDrink).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/status", updateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")

	//Démarrer le serveur
	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(r)); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
