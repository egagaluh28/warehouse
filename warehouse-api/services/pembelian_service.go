package services

import (
    "database/sql"
    "fmt"
    "warehouse-api/models"
    "warehouse-api/repositories"
    "warehouse-api/utils"
)

type PembelianService interface {
    Create(req models.CreatePembelianRequest) (*models.BeliHeader, error)
}

type pembelianService struct {
    db          *sql.DB
    repo        repositories.PembelianRepository
    stokRepo    repositories.StokRepository
    barangRepo  repositories.BarangRepository
}

func NewPembelianService(db *sql.DB, repo repositories.PembelianRepository, stokRepo repositories.StokRepository, barangRepo repositories.BarangRepository) PembelianService {
    return &pembelianService{db, repo, stokRepo, barangRepo}
}

func (s *pembelianService) Create(req models.CreatePembelianRequest) (*models.BeliHeader, error) {
    // 1. Input data pembelian - Validate barang exists
    var totalTrans float64
    var details []models.BeliDetail
    var stockData map[int]int = make(map[int]int) // barangID -> stokSebelum

    for _, d := range req.Details {
        // Validasi: cek apakah barang exists
        barang, err := s.barangRepo.GetByID(d.BarangID)
        if err != nil {
            return nil, fmt.Errorf("barang ID %d tidak ditemukan: %v", d.BarangID, err)
        }
        if barang == nil {
            return nil, fmt.Errorf("barang ID %d tidak ditemukan", d.BarangID)
        }

        // Ambil stok saat ini
        currentStok, _ := s.stokRepo.GetByBarangID(d.BarangID)
        var stokSebelum int
        if currentStok != nil {
            stokSebelum = currentStok.StokAkhir
        } else {
            stokSebelum = 0
        }
        stockData[d.BarangID] = stokSebelum

        // 2. Calculate total
        subtotal := float64(d.Qty) * d.Harga
        totalTrans += subtotal
        
        details = append(details, models.BeliDetail{
            BarangID: d.BarangID,
            Qty:      d.Qty,
            Harga:    d.Harga,
            Subtotal: subtotal,
        })
    }

    // Start transaction
    tx, err := s.db.Begin()
    if err != nil {
        return nil, fmt.Errorf("gagal memulai transaksi database: %v", err)
    }
    defer tx.Rollback()

    // Auto Generate No Faktur
    if req.NoFaktur == "" {
        req.NoFaktur = utils.GenerateNoFakturBeli(s.db)
    }

    header := &models.BeliHeader{
        NoFaktur: req.NoFaktur,
        Supplier: req.Supplier,
        Total:    totalTrans,
        UserID:   req.UserID,
        Status:   "selesai",
    }

    // 3. Update Stok & 4. Record History (SEBELUM save transaction)
    for _, d := range details {
        stokSebelum := stockData[d.BarangID]
        
        // Update stok
        if err := s.stokRepo.CreateOrUpdate(tx, d.BarangID, d.Qty); err != nil {
             return nil, fmt.Errorf("gagal memperbarui stok barang ID %d: %v", d.BarangID, err)
        }
        
        // Hitung stok sesudah
        stokSesudah := stokSebelum + d.Qty
        
        // Record history
        history := &models.HistoryStok{
            BarangID:       d.BarangID,
            UserID:         req.UserID,
            JenisTransaksi: "masuk",
            Jumlah:         d.Qty,
            StokSebelum:    stokSebelum,
            StokSesudah:    stokSesudah,
            Keterangan:     "Pembelian " + header.NoFaktur,
        }
        if err := s.stokRepo.CreateHistory(tx, history); err != nil {
             return nil, fmt.Errorf("gagal membuat riwayat stok: %v", err)
        }
    }

    // 5. Save transaction (beli_header & beli_detail) - PALING AKHIR sebelum commit
    if err := s.repo.Create(tx, header, details); err != nil {
        return nil, fmt.Errorf("gagal membuat transaksi: %v", err)
    }

    // Commit transaksi
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("gagal commit transaksi: %v", err)
    }

    return header, nil
}
