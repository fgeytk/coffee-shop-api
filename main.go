package main

import (
	"coffee-shop-api/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//Routeur mux
	r := mux.NewRouter()

	//Def les routes
	r.HandleFunc("/menu", handlers.GetMenu).Methods("GET")
	r.HandleFunc("/drinks/{id}", handlers.GetDrink).Methods("GET")
	r.HandleFunc("/orders", handlers.CreateOrder).Methods("POST")
	r.HandleFunc("/orders", handlers.GetOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", handlers.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/status", handlers.UpdateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", handlers.DeleteOrder).Methods("DELETE")

	//Démarrer le serveur
	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", handlers.CorsMiddleware(r)); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
