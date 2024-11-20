package models

import (
	"database/sql"
	"fmt"
)

// Product struct represents the product model
type Product struct {
	ID                 int      `json:"id"`
	UserID             int      `json:"user_id"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `json:"product_images"`
	ProductPrice       float64  `json:"product_price"`
}

// Save method saves the product to the database
func (p *Product) Save(db *sql.DB) error {
	query := `INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return db.QueryRow(query, p.UserID, p.ProductName, p.ProductDescription, p.ProductImages, p.ProductPrice).Scan(&p.ID)
}

// GetProductByID fetches a product by its ID
func GetProductByID(db *sql.DB, id string) (*Product, error) {
	var product Product
	query := `SELECT id, user_id, product_name, product_description, product_images, product_price FROM products WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription, &product.ProductImages, &product.ProductPrice)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}
	return &product, nil
}

// GetProducts fetches all products, optionally filtered by user_id, price_min, and price_max
func GetProducts(db *sql.DB, userID string, minPrice, maxPrice float64) ([]Product, error) {
	var products []Product
	query := `SELECT id, user_id, product_name, product_description, product_images, product_price FROM products WHERE 1=1`

	// Add filters to the query
	if userID != "" {
		query += " AND user_id = $1"
	}
	if minPrice > 0 {
		query += " AND product_price >= $2"
	}
	if maxPrice > 0 {
		query += " AND product_price <= $3"
	}

	rows, err := db.Query(query, userID, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription, &product.ProductImages, &product.ProductPrice); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
