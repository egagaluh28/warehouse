package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
    // Load .env file if exists
    err := godotenv.Load()
    if err != nil {
        log.Println("Note: .env file not found or error loading, using system env variables.")
    }

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
    
    // Default values if env not set (Helper for beginner)
    if host == "" { host = "localhost" }
    if port == "" { port = "5432" }
    if user == "" { user = "postgres" }
    if dbname == "" { dbname = "warehouse" }

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v. Please check your DB configuration in .env", err)
	}

	DB = db
	fmt.Println("Successfully connected to the database!")
}
