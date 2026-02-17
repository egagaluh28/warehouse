package seeders

import (
	"database/sql"
	"fmt"
	"log"
)

// SeedBarang populates the database with initial barang data, stock, and transactions
func SeedBarang(db *sql.DB) {
	fmt.Println("Seeding Barang, Stock, and Transactions...")

    // 1. Seed Master Barang
	barangs := []struct {
		KodeBarang string
		NamaBarang string
		Deskripsi  string
		Satuan     string
		HargaBeli  float64
		HargaJual  float64
        StokAwal   int
	}{
		{"BRG001", "Laptop Dell XPS 13", "Laptop Business Grade", "unit", 15000000, 17500000, 10},
		{"BRG002", "Mouse Wireless Logitech", "Mouse Wireless 2.4GHz", "pcs", 250000, 350000, 50},
		{"BRG003", "Keyboard Mechanical", "Keyboard Mechanical RGB", "pcs", 800000, 1200000, 30},
        {"BRG004", "Monitor 24 inch", "Monitor LED 24 inch Full HD", "unit", 2000000, 2800000, 15},
        {"BRG005", "Webcam HD 1080p", "Webcam High Definition", "pcs", 450000, 650000, 25},
	}

	for _, b := range barangs {
		var id int
		err := db.QueryRow("SELECT id FROM master_barang WHERE kode_barang = $1", b.KodeBarang).Scan(&id)

		if err == sql.ErrNoRows {
			// Insert Barang
            var newID int
			err = db.QueryRow("INSERT INTO master_barang (kode_barang, nama_barang, deskripsi, satuan, harga_beli, harga_jual) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
				b.KodeBarang, b.NamaBarang, b.Deskripsi, b.Satuan, b.HargaBeli, b.HargaJual).Scan(&newID)
			
            if err != nil {
				log.Printf("Failed to insert barang %s: %v", b.NamaBarang, err)
                continue
			} 
            
            // Insert Initial Stock
            _, err = db.Exec("INSERT INTO mstok (barang_id, stok_akhir) VALUES ($1, $2)", newID, b.StokAwal)
            if err != nil {
                log.Printf("Failed to insert initial stock for %s: %v", b.NamaBarang, err)
            }

            fmt.Printf("Inserted barang: %s (Stok: %d)\n", b.NamaBarang, b.StokAwal)
		}
	}

    // 2. Seed Pembelian (Beli Header & Detail)
    seedPembelian(db)
    
    // 3. Seed Penjualan (Jual Header & Detail)
    seedPenjualan(db)

    // 4. Seed History Stok
    seedHistory(db)
}

func seedPembelian(db *sql.DB) {
    // Check if data exists first to avoid duplicates
    var exists bool
    db.QueryRow("SELECT EXISTS(SELECT 1 FROM beli_header WHERE no_faktur='BLI001')").Scan(&exists)
    if exists {
        fmt.Println("Pembelian data already seeded.")
        return
    }

    // Insert Beli Header 1
     // Note: Assuming UserID 2 (staff1) and 3 (staff2) exist from user_seeder. Assuming IDs are sequential or fetched.
     // In real seeder we should fetch UserID by username.
    var staff1ID, staff2ID int
    db.QueryRow("SELECT id FROM users WHERE username='staff1'").Scan(&staff1ID)
    db.QueryRow("SELECT id FROM users WHERE username='staff2'").Scan(&staff2ID)

    var beli1ID int
    err := db.QueryRow("INSERT INTO beli_header (no_faktur, supplier, total, user_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        "BLI001", "PT Supplier Elektronik", 32500000, staff1ID, "selesai").Scan(&beli1ID)
    if err == nil {
        // Insert Details
        // Get Barang IDs
        var brg1, brg2 int
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG001'").Scan(&brg1)
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG002'").Scan(&brg2)
        
        db.Exec("INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", beli1ID, brg1, 2, 15000000, 30000000)
        db.Exec("INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", beli1ID, brg2, 10, 250000, 2500000)
    }

    // Insert Beli Header 2
    var beli2ID int
    err = db.QueryRow("INSERT INTO beli_header (no_faktur, supplier, total, user_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        "BLI002", "CV Komputer Jaya", 12500000, staff2ID, "selesai").Scan(&beli2ID)
    if err == nil {
         var brg3, brg4, brg5 int
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG003'").Scan(&brg3)
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG004'").Scan(&brg4)
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG005'").Scan(&brg5)

        db.Exec("INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", beli2ID, brg3, 5, 800000, 4000000)
        db.Exec("INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", beli2ID, brg4, 3, 2000000, 6000000)
        db.Exec("INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", beli2ID, brg5, 4, 450000, 1800000)
    }
     fmt.Println("Pembelian seeded.")
}

func seedPenjualan(db *sql.DB) {
     var exists bool
    db.QueryRow("SELECT EXISTS(SELECT 1 FROM jual_header WHERE no_faktur='JUAL001')").Scan(&exists)
    if exists {
        fmt.Println("Penjualan data already seeded.")
        return
    }

    var staff1ID, staff2ID int
    db.QueryRow("SELECT id FROM users WHERE username='staff1'").Scan(&staff1ID)
    db.QueryRow("SELECT id FROM users WHERE username='staff2'").Scan(&staff2ID)

    // Jual 1
    var jual1ID int
    err := db.QueryRow("INSERT INTO jual_header (no_faktur, customer, total, user_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        "JUAL001", "PT Customer Indonesia", 18700000, staff1ID, "selesai").Scan(&jual1ID)
    if err == nil {
        var brg1, brg2, brg3 int
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG001'").Scan(&brg1)
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG002'").Scan(&brg2)
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG003'").Scan(&brg3)

        db.Exec("INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", jual1ID, brg1, 1, 17500000, 17500000)
        db.Exec("INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", jual1ID, brg2, 2, 350000, 700000)
        db.Exec("INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", jual1ID, brg3, 1, 1200000, 1200000)
    }

    // Jual 2
    var jual2ID int
    err = db.QueryRow("INSERT INTO jual_header (no_faktur, customer, total, user_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        "JUAL002", "CV Tech Solution", 4150000, staff2ID, "selesai").Scan(&jual2ID)
    if err == nil {
         var brg2, brg4 int
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG002'").Scan(&brg2)
        db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG004'").Scan(&brg4)

        db.Exec("INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", jual2ID, brg2, 5, 350000, 1750000)
        db.Exec("INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES ($1, $2, $3, $4, $5)", jual2ID, brg4, 1, 2800000, 2800000)
    }
    fmt.Println("Penjualan seeded.")
}

func seedHistory(db *sql.DB) {
     // Check basic unique constraint logic or just simple check
     // For simplicity, just insert if empty
     var count int
     db.QueryRow("SELECT COUNT(*) FROM history_stok").Scan(&count)
     if count > 0 {
         fmt.Println("History already seeded.")
         return
     }

    var staff1ID, staff2ID int
    db.QueryRow("SELECT id FROM users WHERE username='staff1'").Scan(&staff1ID)
    db.QueryRow("SELECT id FROM users WHERE username='staff2'").Scan(&staff2ID)

    var brg1, brg2, brg3, brg4, brg5 int
    db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG001'").Scan(&brg1)
    db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG002'").Scan(&brg2)
    db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG003'").Scan(&brg3)
    db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG004'").Scan(&brg4)
    db.QueryRow("SELECT id FROM master_barang WHERE kode_barang='BRG005'").Scan(&brg5)

    // Insert History Records manually to match dummy data
    db.Exec("INSERT INTO history_stok (barang_id, user_id, jenis_transaksi, jumlah, stok_sebelum, stok_sesudah, keterangan) VALUES ($1, $2, $3, $4, $5, $6, $7)", brg1, staff1ID, "masuk", 2, 0, 2, "Pembelian BLI001")
    db.Exec("INSERT INTO history_stok (barang_id, user_id, jenis_transaksi, jumlah, stok_sebelum, stok_sesudah, keterangan) VALUES ($1, $2, $3, $4, $5, $6, $7)", brg2, staff1ID, "masuk", 10, 0, 10, "Pembelian BLI001")
    
    // ... Add more as needed
    fmt.Println("History seeded.")
}
