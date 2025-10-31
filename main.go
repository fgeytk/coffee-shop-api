package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

// Drink représente une boisson du menu
type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"` // coffee, tea, cold
	BasePrice float64 `json:"base_price"`
}

// OrderStatus représente l'état d'une commande
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusPickedUp  OrderStatus = "picked-up"
	StatusCancelled OrderStatus = "cancelled"
)

// Order représente une commande
type Order struct {
	ID           string      `json:"id"`
	DrinkID      string      `json:"drink_id"`
	DrinkName    string      `json:"drink_name"`
	Size         string      `json:"size"`   // small, medium, large
	Extras       []string    `json:"extras"` // milk, sugar, cream, caramel
	CustomerName string      `json:"customer_name"`
	Status       OrderStatus `json:"status"`
	TotalPrice   float64     `json:"total_price"`
	OrderedAt    time.Time   `json:"ordered_at"`
}

// Base de données en mémoire
var drinks = []Drink{
	{ID: "1", Name: "Espresso", Category: "coffee", BasePrice: 2.0},
	{ID: "2", Name: "Cappuccino", Category: "coffee", BasePrice: 3.0},
	{ID: "3", Name: "Latte", Category: "coffee", BasePrice: 3.5},
	{ID: "4", Name: "Black Tea", Category: "tea", BasePrice: 2.5},
	{ID: "5", Name: "Green Tea", Category: "tea", BasePrice: 2.5},
	{ID: "6", Name: "Iced Coffee", Category: "cold", BasePrice: 3.0},
	{ID: "7", Name: "Iced Tea", Category: "cold", BasePrice: 2.5},
}
var orders []Order
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
	json.NewEncoder(w).Encode(drinks)
	fmt.Println("Menu :")
}

func getDrink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	drinkID := vars["id"]
	for _, drink := range drinks {
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
	var newOrder Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	//gerer l'erreur
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Données de commande invalides"})
		return
	}
	//verifie si la boisson existe
	var foundDrink *Drink
	for _, drink := range drinks {
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
	newOrder.Status = StatusPending
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
		Status OrderStatus `json:"status"`
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
			if order.Status == StatusPickedUp {
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
