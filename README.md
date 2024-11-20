# Product Management System

This repository contains the code for a **Product Management System**. The system allows users to manage products, including creating new products, retrieving product information, and querying all products. Additionally, the system integrates image processing, caching with Redis, and database interactions.

## Features

- **Create Product**: Allows users to create a new product.
- **Get Product**: Fetches product details by product ID.
- **Get All Products**: Retrieves a list of all products with optional filters (e.g., user ID, price range).
- **Image Processing**: Simulated image processing using an external image processor service.
- **Redis Caching**: Caches product data to reduce database load and improve performance.

## Project Structure

Here's a brief overview of the directory structure:

```
product-management/
├── api/                        # Contains API-related code (handlers, routes, middleware)
│   ├── handlers/               # HTTP handlers for API endpoints
│   ├── middleware/             # Custom middleware
│   └── routes.go               # API route definitions
├── db/                         # Database-related code
│   ├── migrations/             # SQL migration files
│   └── connection.go           # Database connection setup
├── image-processor/            # Code related to image processing
├── cache/                      # Redis caching code
├── config/                     # Configuration files (env variables, config loading)
├── models/                     # Data models (Product, User, etc.)
├── tests/                      # Unit and integration tests
├── logs/                       # Logs for the application
├── main.go                     # Main entry point for the application
├── go.mod                      # Go module dependencies
└── go.sum                      # Go module checksum
```

## Prerequisites

Before setting up the system, ensure you have the following installed:

- **Go (1.18+)**: The programming language used for this project. [Install Go](https://golang.org/doc/install).
- **PostgreSQL**: The database used to store product information. [Install PostgreSQL](https://www.postgresql.org/download/).
- **Redis**: Used for caching product data. [Install Redis](https://redis.io/download).
- **Docker** (optional but recommended): For running services in containers (PostgreSQL, Redis).

## Setup Instructions

### 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/kavin-t28/product-management.git
cd product-management
```

### 2. Set Up the Environment

Create a `.env` file in the root of the project directory. This file will hold the environment variables required by the application:

```bash
cp .env.example .env
```

Modify the `.env` file with your own values for the following variables:

- `DB_HOST`: Hostname of your PostgreSQL database (e.g., `localhost`).
- `DB_PORT`: Port for the PostgreSQL database (default is `5432`).
- `DB_NAME`: Database name (e.g., `product_management`).
- `DB_USER`: Database user (e.g., `pm_user`).
- `DB_PASSWORD`: Password for the database user.
- `REDIS_HOST`: Hostname of your Redis server (e.g., `localhost`).
- `REDIS_PORT`: Port for Redis (default is `6379`).

Example `.env` file:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=product_management
DB_USER=pm_user
DB_PASSWORD=yourpassword
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 3. Set Up PostgreSQL Database

Make sure PostgreSQL is installed and running on your machine. You can create the required database and user as follows:

```bash
# Access PostgreSQL shell
psql -U postgres

# Create database
CREATE DATABASE product_management;

# Create user
CREATE USER pm_user WITH PASSWORD 'yourpassword';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE product_management TO pm_user;

# Exit PostgreSQL
\q
```

### 4. Apply Database Migrations

Use the SQL migration files in the `db/migrations/` folder to set up the database schema. Run the migrations using the following commands:

```bash
psql -U pm_user -d product_management -f db/migrations/001_create_users_table.sql
psql -U pm_user -d product_management -f db/migrations/002_create_products_table.sql
```

### 5. Install Dependencies

Install the required Go dependencies:

```bash
go mod tidy
```

### 6. Running the Application

Once the database and Redis server are set up and the environment is configured, you can run the application:

```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

### 7. Running Tests

To run the tests, you can use the `go test` command. You can run all tests or specific test files:

```bash
# Run all tests
go test ./...

# Run a specific test file
go test ./tests/api_test.go
```

## API Endpoints

Here are the available API endpoints:

### 1. `POST /products`
Create a new product.

#### Request body:
```json
{
  "user_id": 1,
  "product_name": "Product Name",
  "product_description": "Description of the product",
  "product_images": ["image1.jpg", "image2.jpg"],
  "product_price": 19.99
}
```

#### Response:
```json
{
  "id": 1,
  "user_id": 1,
  "product_name": "Product Name",
  "product_description": "Description of the product",
  "product_images": ["image1.jpg", "image2.jpg"],
  "product_price": 19.99
}
```

### 2. `GET /products/{id}`
Get product details by product ID.

#### Response:
```json
{
  "id": 1,
  "user_id": 1,
  "product_name": "Product Name",
  "product_description": "Description of the product",
  "product_images": ["image1.jpg", "image2.jpg"],
  "product_price": 19.99
}
```

### 3. `GET /products`
Get all products with optional query parameters for filtering.

#### Query parameters:
- `user_id`: Filter by user ID.
- `price_min`: Filter by minimum price.
- `price_max`: Filter by maximum price.

Example request:
```
GET /products?user_id=1&price_min=10&price_max=100
```

#### Response:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "product_name": "Product Name",
    "product_description": "Description of the product",
    "product_images": ["image1.jpg", "image2.jpg"],
    "product_price": 19.99
  }
]
```

## System Architecture

### 1. **Product Model**: 
The product model is a simple struct with fields like `ID`, `UserID`, `ProductName`, `ProductDescription`, `ProductImages`, and `ProductPrice`.

### 2. **Database**:
The PostgreSQL database stores product data. We use the `products` table to store product information and related details. The database is connected via the `db/connection.go` file.

### 3. **Redis Cache**:
We use Redis as a caching layer to store product data. When a product is fetched by ID, we first check if it exists in the Redis cache. If not, we query the database and store the result in Redis for future requests.

### 4. **Image Processing**:
When a new product is created, the system simulates image processing (e.g., compression or uploading to an image storage service) via the `image-processor/queue.go` and `image-processor/processor.go`.

### 5. **Middleware**:
Custom middleware like logging is added in the `api/middleware/logging.go` file.



## Troubleshooting

### 1. **Redis Connection Issues**
   - **Error**: `redis: dial tcp: lookup localhost: no such host`
   - **Cause**: The Redis server is not running or the connection details in the `.env` file are incorrect.
   - **Solution**:
     - Ensure Redis is running on the specified `REDIS_HOST` and `REDIS_PORT` (default is `localhost:6379`).
     - Use the `redis-cli` tool to test if you can connect to Redis:
       ```bash
       redis-cli -h localhost -p 6379
       ```
     - If Redis is running in Docker, ensure that the Docker container's network is accessible from your Go application.

### 2. **Database Connection Issues**
   - **Error**: `pq: connection to server failed`
   - **Cause**: PostgreSQL server is down, misconfigured, or connection details in `.env` are incorrect.
   - **Solution**:
     - Check if PostgreSQL is running:
       ```bash
       sudo systemctl status postgresql
       ```
     - Verify the connection details in your `.env` file (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, etc.).
     - If using Docker for PostgreSQL, ensure the container is running and the correct ports are exposed.
     - Try connecting manually to PostgreSQL to check the credentials:
       ```bash
       psql -U pm_user -d product_management -h localhost
       ```

### 3. **Database Migrations Not Working**
   - **Error**: `psql: command not found` or `ERROR: permission denied for schema public`
   - **Cause**: You might not have applied the database migrations correctly, or you might lack sufficient permissions.
   - **Solution**:
     - Ensure you have the necessary permissions for the PostgreSQL user (`pm_user`) to create tables.
     - Verify that you have run all migrations using the `psql` commands.
     - If running into permission errors, check the PostgreSQL user roles and ensure `pm_user` has `CREATE` and `INSERT` permissions.
     - Re-run the migration SQL scripts manually:
       ```bash
       psql -U pm_user -d product_management -f db/migrations/001_create_users_table.sql
       psql -U pm_user -d product_management -f db/migrations/002_create_products_table.sql
       ```

### 4. **Image Processing Not Working**
   - **Error**: `Error processing image: network timeout` or `unable to access external image URL`
   - **Cause**: The image processing service cannot reach the image URL or an external service due to network issues.
   - **Solution**:
     - Check if the image URL is accessible by trying to download the image manually.
     - Ensure that any external image processing or cloud services (e.g., AWS S3) are configured correctly.
     - If using a simulated image processor, verify that the path is correctly defined in the `image-processor/processor.go` file and that no network issues are affecting the simulation.

### 5. **API Route Not Found**
   - **Error**: `404 page not found`
   - **Cause**: The API endpoint you're trying to access doesn't exist or hasn't been properly registered in the router.
   - **Solution**:
     - Double-check the route definitions in `api/routes.go` and ensure that the correct HTTP methods (`POST`, `GET`) are used.
     - Verify that the handler functions (`CreateProductHandler`, `GetProductHandler`, etc.) are correctly registered with the router.
     - Use a tool like `curl` or Postman to check if the endpoint is correctly exposed:
       ```bash
       curl -X GET http://localhost:8080/products
       ```

### 6. **Cache Not Being Used**
   - **Error**: Slow response times for product data fetching even though Redis is enabled.
   - **Cause**: The Redis cache may not be used correctly, or the cache may be expired.
   - **Solution**:
     - Ensure that the caching logic is properly implemented in `services/cache.go`. When fetching a product, check if it's correctly retrieving the cached data.
     - If the cache isn't found, ensure that the logic correctly populates the cache after fetching data from the database.
     - You can check the cache status using the `redis-cli` tool:
       ```bash
       redis-cli GET product:1
       ```
     - Make sure you have added cache expiry or invalidation logic to avoid serving stale data.

### 7. **CORS Issues (Cross-Origin Resource Sharing)**
   - **Error**: `No 'Access-Control-Allow-Origin' header is present on the requested resource`
   - **Cause**: The frontend is trying to access the API, but CORS is not enabled on the server.
   - **Solution**:
     - Add middleware for CORS handling in the `api/middleware/cors.go` (or equivalent). You can use the `github.com/rs/cors` package for this.
     - Example of enabling CORS:
       ```go
       package middleware

       import (
           "github.com/rs/cors"
           "net/http"
       )

       func EnableCORS(next http.Handler) http.Handler {
           c := cors.New(cors.Options{
               AllowedOrigins: []string{"*"},
               AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
               AllowedHeaders: []string{"Content-Type", "Authorization"},
           })
           return c.Handler(next)
       }
       ```
     - Apply this middleware in the router:
       ```go
       router.Use(middleware.EnableCORS)
       ```

### 8. **Missing or Incorrect Environment Variables**
   - **Error**: `panic: missing environment variable DB_HOST`
   - **Cause**: Required environment variables (such as database or Redis connection details) are not defined or are incorrect.
   - **Solution**:
     - Ensure that the `.env` file is present and correctly configured.
     - Run the following command to ensure all environment variables are loaded:
       ```bash
       source .env
       ```
     - Double-check the contents of `.env` to ensure all required variables are set.

### 9. **Go Modules or Dependency Issues**
   - **Error**: `go: module github.com/some/module@v0.0.0: cannot find module providing package`
   - **Cause**: A missing or outdated Go module dependency.
   - **Solution**:
     - Run `go mod tidy` to clean up any outdated or unused dependencies.
     - If the error persists, try running `go get` to fetch the missing dependencies:
       ```bash
       go get github.com/some/module
       ```
     - Ensure that the `go.mod` and `go.sum` files are committed and up-to-date.

### 10. **Port Already in Use**
   - **Error**: `Error: listen tcp :8080: bind: address already in use`
   - **Cause**: Another application is already running on the same port.
   - **Solution**:
     - Check for processes using port 8080:
       ```bash
       sudo lsof -i :8080
       ```
     - Kill the process occupying the port:
       ```bash
       kill <PID>
       ```
     - Alternatively, change the port number in the `main.go` file or use a different port when starting the server:
       ```go
       log.Fatal(http.ListenAndServe(":8081", router))
       ```

### 11. **Incorrect Product Data or Missing Fields**
   - **Error**: `Invalid product data`, `Product not found`
   - **Cause**: The product data passed in the request body or database query is malformed or incomplete.
   - **Solution**:
     - Ensure the request body matches the expected format, including all required fields (e.g., `user_id`, `product_name`, `product_price`).
     - Check if the product exists in the database by running a direct SQL query:
       ```sql
       SELECT * FROM products WHERE id = <product_id>;
       ```
     - If the product is not found, ensure the ID is correct and that the product was successfully created.

---
