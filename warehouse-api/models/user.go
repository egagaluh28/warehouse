package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"newstaff"`
	Password string `json:"password" binding:"required" example:"staff123"`
	Email    string `json:"email" binding:"required,email" example:"staff@example.com"`
	FullName string `json:"full_name" binding:"required" example:"New Staff Member"`
	Role     string `json:"role" binding:"required" example:"staff"` // admin or staff
}
