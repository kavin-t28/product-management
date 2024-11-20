package main

import (
	"log"
	"net/http"
	"product-management/api"
	"product-management/cache"
	"product-management/config"
	"product-management/db"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database connection
	db.InitDB()

	// Initialize Redis client
	cache.InitRedis()

	// Set up router and routes
	router := mux.NewRouter()
	api.RegisterRoutes(router)

	// Start the server
	log.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
