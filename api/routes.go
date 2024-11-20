package api

import (
	"product-management/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
}
