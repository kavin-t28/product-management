package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-management/api/handlers"
	"product-management/models"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	// Define a test product
	product := models.Product{
		UserID:             1,
		ProductName:        "Test Product",
		ProductDescription: "A description of the test product",
		ProductImages:      []string{"http://example.com/image1.jpg"},
		ProductPrice:       100.00,
	}

	// Convert product struct to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling product: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(productJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a new ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a mock handler for the CreateProduct function
	handler := http.HandlerFunc(handlers.CreateProduct)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response code
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %v, but got %v", http.StatusCreated, rr.Code)
	}

	// Check if the returned product contains the correct data
	var returnedProduct models.Product
	err = json.NewDecoder(rr.Body).Decode(&returnedProduct)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	if returnedProduct.ProductName != product.ProductName {
		t.Errorf("Expected product name %v, but got %v", product.ProductName, returnedProduct.ProductName)
	}
}
func TestGetProduct(t *testing.T) {
	// Mock product data for the test
	product := models.Product{
		ID:                 1,
		UserID:             1,
		ProductName:        "Test Product",
		ProductDescription: "Test product description",
		ProductImages:      []string{"http://example.com/image1.jpg"},
		ProductPrice:       100.00,
	}

	// Mock database query to return the test product (You can mock db.DB.QueryRow here)
	// Here, we'll just assume this product is returned from the database

	// Create a new HTTP request for GET /products/{id}
	req := httptest.NewRequest("GET", "/products/1", nil)

	// Create a new ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Mock the handler for GetProduct
	handler := http.HandlerFunc(handlers.GetProduct)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, rr.Code)
	}

	// Check the returned product data
	var returnedProduct models.Product
	err := json.NewDecoder(rr.Body).Decode(&returnedProduct)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	if returnedProduct.ProductName != product.ProductName {
		t.Errorf("Expected product name %v, but got %v", product.ProductName, returnedProduct.ProductName)
	}
}
func TestGetAllProducts(t *testing.T) {
	// Mock product data for the test
	products := []models.Product{
		{
			ID:                 1,
			UserID:             1,
			ProductName:        "Test Product 1",
			ProductDescription: "Description for test product 1",
			ProductImages:      []string{"http://example.com/image1.jpg"},
			ProductPrice:       50.00,
		},
		{
			ID:                 2,
			UserID:             1,
			ProductName:        "Test Product 2",
			ProductDescription: "Description for test product 2",
			ProductImages:      []string{"http://example.com/image2.jpg"},
			ProductPrice:       150.00,
		},
	}

	// Mock database query to return the list of products (You can mock db.DB.Query here)
	// In this case, we'll assume the mock query would return the above products

	// Create a new HTTP request for GET /products
	req := httptest.NewRequest("GET", "/products", nil)

	// Create a new ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Mock the handler for GetAllProducts
	handler := http.HandlerFunc(handlers.GetAllProducts)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, rr.Code)
	}

	// Check the returned products data
	var returnedProducts []models.Product
	err := json.NewDecoder(rr.Body).Decode(&returnedProducts)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Validate if the returned products match the mocked products
	if len(returnedProducts) != len(products) {
		t.Errorf("Expected %d products, but got %d", len(products), len(returnedProducts))
	}

	for i, product := range products {
		if returnedProducts[i].ProductName != product.ProductName {
			t.Errorf("Expected product name %v, but got %v", product.ProductName, returnedProducts[i].ProductName)
		}
	}
}
