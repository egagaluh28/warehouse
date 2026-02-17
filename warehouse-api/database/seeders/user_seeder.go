package seeders

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedUsers populates the database with initial users
func SeedUsers(db *sql.DB) {
	fmt.Println("Seeding Users...")

	users := []struct {
		Username string
		Password string
		Email    string
		FullName string
		Role     string
	}{
		{"admin", "admin", "admin@warehouse.com", "Administrator System", "admin"},
		{"staff1", "staff", "staff1@warehouse.com", "Staff Gudang A", "staff"},
		{"staff2", "staff", "staff2@warehouse.com", "Staff Gudang B", "staff"},
	}

	for _, u := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.Username, err)
			continue
		}

		// Check if user exists
		var id int
		err = db.QueryRow("SELECT id FROM users WHERE username = $1", u.Username).Scan(&id)

		if err == sql.ErrNoRows {
			// Insert
			_, err = db.Exec("INSERT INTO users (username, password, email, full_name, role) VALUES ($1, $2, $3, $4, $5)",
				u.Username, string(hashedPassword), u.Email, u.FullName, u.Role)
			if err != nil {
				log.Printf("Failed to insert user %s: %v", u.Username, err)
			} else {
				fmt.Printf("Inserted user: %s (Password: %s)\n", u.Username, u.Password)
			}
		} else if err == nil {
			// Update (optional, but good for resetting state)
            // Uncomment to force update on seed
			_, err = db.Exec("UPDATE users SET password = $1, email = $2, full_name = $3, role = $4 WHERE username = $5",
				string(hashedPassword), u.Email, u.FullName, u.Role, u.Username)
             fmt.Printf("User %s updated with new password: %s\n", u.Username, u.Password)
		} else {
			log.Printf("Error checking user %s: %v", u.Username, err)
		}
	}
}
