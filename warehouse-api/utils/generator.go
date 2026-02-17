package utils

import (
    "database/sql"
    "fmt"
    "time"
    "math/rand"
)

// RandomString generates a random string of specified length
func RandomString(length int) string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

// GenerateCode generates a generic code with prefix
// Format: PREFIX-YYMMDD-RANDOM
func GenerateCode(prefix string) string {
    timestamp := time.Now().Format("060102") // YYMMDD
    random := RandomString(4)
    return fmt.Sprintf("%s-%s-%s", prefix, timestamp, random)
}

// GenerateKodeBarang generates a code like BRG-YYMMDD-RANDOM
func GenerateKodeBarang(db *sql.DB) string {
    return GenerateCode("BRG")
}

// GenerateNoFakturBeli generates a code like BELI-YYMMDD-RANDOM
func GenerateNoFakturBeli(db *sql.DB) string {
    return GenerateCode("BELI")
}

// GenerateNoFakturJual generates a code like JUAL-YYMMDD-RANDOM
func GenerateNoFakturJual(db *sql.DB) string {
    return GenerateCode("JUAL")
}
