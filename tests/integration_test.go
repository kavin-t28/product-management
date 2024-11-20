package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"product-management/api/handlers"
	"product-management/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Start by testing the full flow: Create Product -> Get Product -> Get All Products

	// Step 1: Create a product
	product := models.Product{
		UserID:             1,
		ProductName:        "Test Product",
		ProductDescription: "This is a test product for integration testing.",
		ProductImages:      []string{"http://example.com/image.jpg"},
		ProductPrice:       50.00,
	}

	// Convert product struct to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshalling product: %v", err)
	}

	// Create a new HTTP request for POST /products
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(productJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a new ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create the mock handler for the CreateProduct function
	handler := http.HandlerFunc(handlers.CreateProduct)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Step 2: Validate that the product was created successfully
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status %v, but got %v", http.StatusCreated, rr.Code)

	var createdProduct models.Product
	err = json.NewDecoder(rr.Body).Decode(&createdProduct)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Make sure the product data returned is correct
	assert.Equal(t, product.ProductName, createdProduct.ProductName)
	assert.Equal(t, product.ProductDescription, createdProduct.ProductDescription)

	// Step 3: Now test getting the created product using its ID
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/products/%d", createdProduct.ID), nil)
	getRR := httptest.NewRecorder()

	// Handler for GetProduct
	getHandler := http.HandlerFunc(handlers.GetProduct)
	getHandler.ServeHTTP(getRR, getReq)

	// Validate the GetProduct response
	assert.Equal(t, http.StatusOK, getRR.Code, "Expected status %v, but got %v", http.StatusOK, getRR.Code)

	var fetchedProduct models.Product
	err = json.NewDecoder(getRR.Body).Decode(&fetchedProduct)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check if the fetched product matches the created one
	assert.Equal(t, createdProduct.ProductName, fetchedProduct.ProductName)
	assert.Equal(t, createdProduct.ProductDescription, fetchedProduct.ProductDescription)

	// Step 4: Test fetching all products to make sure the newly created product is included
	allProductsReq := httptest.NewRequest("GET", "/products", nil)
	allProductsRR := httptest.NewRecorder()

	// Handler for GetAllProducts
	allProductsHandler := http.HandlerFunc(handlers.GetAllProducts)
	allProductsHandler.ServeHTTP(allProductsRR, allProductsReq)

	// Validate the GetAllProducts response
	assert.Equal(t, http.StatusOK, allProductsRR.Code, "Expected status %v, but got %v", http.StatusOK, allProductsRR.Code)

	var allProducts []models.Product
	err = json.NewDecoder(allProductsRR.Body).Decode(&allProducts)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check that the created product is included in the list of all products
	assert.Len(t, allProducts, 1, "Expected 1 product, but got %v", len(allProducts))
	assert.Equal(t, createdProduct.ProductName, allProducts[0].ProductName)
}
