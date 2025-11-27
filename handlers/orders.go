package handlers

import (
	"coffee-shop-api/database"
	"coffee-shop-api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Données de commande invalides"})
		return
	}

	// Vérifier si la boisson existe dans MySQL
	var drink models.Drink
	err = database.DB.QueryRow("SELECT id, name, base_price FROM drinks WHERE id = ?", newOrder.DrinkID).
		Scan(&drink.ID, &drink.Name, &drink.BasePrice)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Boisson introuvable"})
		return
	}

	newOrder.DrinkName = drink.Name
	newOrder.Status = models.StatusPending
	newOrder.OrderedAt = time.Now()
	newOrder.TotalPrice = calculatePrice(drink.BasePrice, newOrder.Size, newOrder.Extras)

	// Insérer la commande dans MySQL
	extrasJSON := strings.Join(newOrder.Extras, ",")
	result, err := database.DB.Exec(
		"INSERT INTO orders (drink_id, drink_name, size, extras, customer_name, status, total_price, ordered_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		newOrder.DrinkID, newOrder.DrinkName, newOrder.Size, extrasJSON, newOrder.CustomerName, newOrder.Status, newOrder.TotalPrice, newOrder.OrderedAt,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la création de la commande"})
		return
	}

	id, _ := result.LastInsertId()
	newOrder.ID = fmt.Sprintf("%d", id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
	fmt.Printf("Nouvelle commande créée dans MySQL: %+v\n", newOrder)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT id, drink_id, drink_name, size, extras, customer_name, status, total_price, ordered_at FROM orders")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la récupération des commandes"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var extrasStr string
		err := rows.Scan(&order.ID, &order.DrinkID, &order.DrinkName, &order.Size, &extrasStr, &order.CustomerName, &order.Status, &order.TotalPrice, &order.OrderedAt)
		if err != nil {
			continue
		}
		if extrasStr != "" {
			order.Extras = strings.Split(extrasStr, ",")
		}
		orders = append(orders, order)
	}

	json.NewEncoder(w).Encode(orders)
	fmt.Println("Liste des commandes récupérée depuis MySQL")
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]

	var order models.Order
	var extrasStr string
	err := database.DB.QueryRow("SELECT id, drink_id, drink_name, size, extras, customer_name, status, total_price, ordered_at FROM orders WHERE id = ?", orderId).
		Scan(&order.ID, &order.DrinkID, &order.DrinkName, &order.Size, &extrasStr, &order.CustomerName, &order.Status, &order.TotalPrice, &order.OrderedAt)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la récupération de la commande"})
		return
	}

	if extrasStr != "" {
		order.Extras = strings.Split(extrasStr, ",")
	}

	json.NewEncoder(w).Encode(order)
	fmt.Printf("Détails de la commande %s récupérés depuis MySQL\n", orderId)
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

	result, err := database.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", statusUpdate.Status, orderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la mise à jour"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
		return
	}

	// Récupérer la commande mise à jour
	var order models.Order
	var extrasStr string
	database.DB.QueryRow("SELECT id, drink_id, drink_name, size, extras, customer_name, status, total_price, ordered_at FROM orders WHERE id = ?", orderId).
		Scan(&order.ID, &order.DrinkID, &order.DrinkName, &order.Size, &extrasStr, &order.CustomerName, &order.Status, &order.TotalPrice, &order.OrderedAt)

	if extrasStr != "" {
		order.Extras = strings.Split(extrasStr, ",")
	}

	json.NewEncoder(w).Encode(order)
	fmt.Printf("Mise à jour de la commande %s dans MySQL\n", orderId)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderId := mux.Vars(r)["id"]

	// Vérifier le statut de la commande
	var status models.OrderStatus
	err := database.DB.QueryRow("SELECT status FROM orders WHERE id = ?", orderId).Scan(&status)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Commande introuvable"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la vérification"})
		return
	}

	if status == models.StatusPickedUp {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Impossible d'annuler une commande déjà récupérée"})
		return
	}

	// Supprimer la commande
	_, err = database.DB.Exec("DELETE FROM orders WHERE id = ?", orderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la suppression"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Printf("Commande %s supprimée de MySQL\n", orderId)
}
