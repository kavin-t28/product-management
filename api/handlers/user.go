package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"product-management/db"
	"product-management/models"

	"github.com/gorilla/mux"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price) 
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = db.DB.QueryRow(query, product.UserID, product.ProductName, product.ProductDescription, product.ProductImages, product.ProductPrice).Scan(&product.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var product models.Product
	query := `SELECT id, user_id, product_name, product_description, product_images, product_price FROM products WHERE id = $1`
	err := db.DB.QueryRow(query, id).Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription, &product.ProductImages, &product.ProductPrice)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, user_id, product_name, product_description, product_images, product_price FROM products")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching products: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription, &product.ProductImages, &product.ProductPrice); err != nil {
			http.Error(w, fmt.Sprintf("Error reading product data: %v", err), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
