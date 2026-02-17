package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
    "os"
)

func main() {
    password := "admin123"
    if len(os.Args) > 1 {
        password = os.Args[1]
    }

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    hash := string(bytes)
	fmt.Println("Password:", password)
    fmt.Println("Hash:", hash)
}
