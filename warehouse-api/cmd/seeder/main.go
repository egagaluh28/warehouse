package main

import (
	"flag"
	"fmt"
	"warehouse-api/config"
	"warehouse-api/database/seeders"
)

func main() {
    // Define flags
    runSeed := flag.Bool("seed", false, "Run database seeders")
    flag.Parse()

	// Connect to Database
	config.ConnectDB()
    defer config.DB.Close()

    if *runSeed {
        runSeeders()
        return
    }

    fmt.Println("Use -seed flag to run seeders")
}

func runSeeders() {
    fmt.Println("--- Starting Database Seeding ---")
    
    seeders.SeedUsers(config.DB)
    seeders.SeedBarang(config.DB)
    
    fmt.Println("--- Database Seeding Completed ---")
}
