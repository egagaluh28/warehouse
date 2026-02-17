package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"warehouse-api/handlers"
	"warehouse-api/models"
	"warehouse-api/repositories"
	"warehouse-api/services"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5432")
	dbUser := getEnv("TEST_DB_USER", "postgres")
	dbPass := getEnv("TEST_DB_PASSWORD", "admin")
	dbName := getEnv("TEST_DB_NAME", "warehouse_test")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Skipping integration tests")
		os.Exit(0)
	}

	if err := testDB.Ping(); err != nil {
		fmt.Println("Skipping integration tests")
		os.Exit(0)
	}

	setupTestSchema(testDB)

	code := m.Run()

	cleanupTestSchema(testDB)
	testDB.Close()

	os.Exit(code)
}

func setupTestSchema(db *sql.DB) {
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100),
			full_name VARCHAR(100),
			role VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS barang (
			id SERIAL PRIMARY KEY,
			kode_barang VARCHAR(50) UNIQUE NOT NULL,
			nama_barang VARCHAR(100) NOT NULL,
			satuan VARCHAR(20) NOT NULL,
			harga_beli DECIMAL(15,2) NOT NULL,
			harga_jual DECIMAL(15,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
}

func cleanupTestSchema(db *sql.DB) {
	_, _ = db.Exec("TRUNCATE users, barang CASCADE")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func TestDatabaseConnection(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	t.Run("Ping database", func(t *testing.T) {
		assert.NoError(t, testDB.Ping())
	})
}

func TestAuthIntegration(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	testDB.Exec("TRUNCATE users CASCADE")

	userRepo := repositories.NewUserRepository(testDB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	registerReq := models.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
		FullName: "Test User",
		Role:     "staff",
	}

	user, err := userService.Register(&registerReq)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	loginReq := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
