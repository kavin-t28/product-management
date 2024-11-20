package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	models "product-management/services"
	"product-management/utils"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// RegisterProductHandlers sets up the routes for the product API
func RegisterProductHandlers(router *mux.Router, db *sql.DB, cache *redis.Client) {
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		CreateProductHandler(w, r, db, cache)
	}).Methods("POST")

	router.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetProductHandler(w, r, db, cache)
	}).Methods("GET")

	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		GetProductsHandler(w, r, db)
	}).Methods("GET")
}

// CreateProductHandler handles product creation
func CreateProductHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, cache *redis.Client) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Save product to DB
	if err := product.Save(db); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save product")
		return
	}

	// Queue image processing
	// services.QueueImageProcessing(product.ProductImages)

	utils.RespondWithJSON(w, http.StatusCreated, product)
}

// GetProductHandler fetches a product by ID
func GetProductHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, cache *redis.Client) {
	id := mux.Vars(r)["id"]

	// Fetch from DB if not found in cache
	product, err := models.GetProductByID(db, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Cache the result for future use
	// services.CacheProduct(cache, product)

	utils.RespondWithJSON(w, http.StatusOK, product)
}

// GetProductsHandler retrieves all products with optional filtering
func GetProductsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.URL.Query().Get("user_id")
	priceMin := r.URL.Query().Get("price_min")
	priceMax := r.URL.Query().Get("price_max")

	// Parse the price filters as floats if present
	var minPrice, maxPrice float64
	var err error

	if priceMin != "" {
		minPrice, err = utils.ParsePrice(priceMin)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid price_min filter")
			return
		}
	}

	if priceMax != "" {
		maxPrice, err = utils.ParsePrice(priceMax)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid price_max filter")
			return
		}
	}

	// Get products from DB
	products, err := models.GetProducts(db, userID, minPrice, maxPrice)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve products: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}
