package db

import (
	"database/sql"
	"fmt"
	"log"
	"product-management/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.DBUser, config.DBPassword, config.DBName, config.DBHost, config.DBPort)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	log.Println("Connected to the database")
}
