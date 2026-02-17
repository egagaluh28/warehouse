package services

import (
    "database/sql"
    "fmt"
    "warehouse-api/models"
    "warehouse-api/repositories"
    "warehouse-api/utils"
)

type PenjualanService interface {
    Create(req models.CreatePenjualanRequest) (*models.JualHeader, error)
}

type penjualanService struct {
    db          *sql.DB
    repo        repositories.PenjualanRepository
    stokRepo    repositories.StokRepository
    barangRepo  repositories.BarangRepository
}

func NewPenjualanService(db *sql.DB, repo repositories.PenjualanRepository, stokRepo repositories.StokRepository, barangRepo repositories.BarangRepository) PenjualanService {
    return &penjualanService{db, repo, stokRepo, barangRepo}
}

func (s *penjualanService) Create(req models.CreatePenjualanRequest) (*models.JualHeader, error) {
    // 1. Input data penjualan - Validate barang exists & Check stock availability
    var totalTrans float64
    var details []models.JualDetail
    var stockData map[int]int = make(map[int]int) // barangID -> stokSebelum

    for _, d := range req.Details {
        // Validasi: cek apakah barang exists dan dapatkan stok
        currentStok, err := s.stokRepo.GetByBarangID(d.BarangID)
        if err != nil {
            return nil, fmt.Errorf("barang ID %d tidak ditemukan dalam stok", d.BarangID)
        }
        if currentStok == nil {
            return nil, fmt.Errorf("barang ID %d tidak ditemukan dalam stok", d.BarangID)
        }
        
        // Check stock availability
        if currentStok.StokAkhir < d.Qty {
            return nil, fmt.Errorf("stok tidak mencukupi untuk Barang ID %d. Tersedia: %d, Diminta: %d", d.BarangID, currentStok.StokAkhir, d.Qty)
        }

        stokSebelum := currentStok.StokAkhir
        stockData[d.BarangID] = stokSebelum

        // 2. Calculate total
        subtotal := float64(d.Qty) * d.Harga
        totalTrans += subtotal
        
        details = append(details, models.JualDetail{
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
        req.NoFaktur = utils.GenerateNoFakturJual(s.db)
    }

    header := &models.JualHeader{
        NoFaktur: req.NoFaktur, 
        Customer: req.Customer,
        Total:    totalTrans,
        UserID:   req.UserID,
        Status:   "selesai",
    }

    // 3. Update Stok & 4. Record History (SEBELUM save transaction)
    for _, d := range details {
        stokSebelum := stockData[d.BarangID]
        
        // Update stok (kurangi)
        if err := s.stokRepo.CreateOrUpdate(tx, d.BarangID, -d.Qty); err != nil {
             return nil, fmt.Errorf("gagal memperbarui stok barang ID %d: %v", d.BarangID, err)
        }
        
        // Hitung stok sesudah
        stokSesudah := stokSebelum - d.Qty
        
        // Record history
        history := &models.HistoryStok{
            BarangID:       d.BarangID,
            UserID:         req.UserID,
            JenisTransaksi: "keluar",
            Jumlah:         d.Qty,
            StokSebelum:    stokSebelum,
            StokSesudah:    stokSesudah,
            Keterangan:     "Penjualan " + header.NoFaktur,
        }
        if err := s.stokRepo.CreateHistory(tx, history); err != nil {
             return nil, fmt.Errorf("gagal membuat riwayat stok: %v", err)
        }
    }

    // 5. Save transaction (jual_header & jual_detail) - PALING AKHIR sebelum commit
    if err := s.repo.Create(tx, header, details); err != nil {
        return nil, fmt.Errorf("gagal membuat transaksi: %v", err)
    }

    // Commit transaksi
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("gagal commit transaksi: %v", err)
    }

    return header, nil
}
