package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
    "os"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"
const RoleKey contextKey = "role"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
            // Header Authorization diperlukan
			http.Error(w, "Otorisasi diperlukan", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        secret := []byte(os.Getenv("JWT_SECRET"))
        if len(secret) == 0 {
            secret = []byte("supersecretkey") // Fallback
        }

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
            // Token tidak valid
			http.Error(w, "Token tidak valid", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
            // Klaim token tidak valid
			http.Error(w, "Klaim token tidak valid", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
             http.Error(w, "User ID tidak valid dalam token", http.StatusUnauthorized)
             return
        }
        userID := int(userIDFloat)
        
		role, ok := claims["role"].(string)
        if !ok {
             role = "staff"
        }

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
