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
    tx, err := s.db.Begin()
    if err != nil {
        return nil, fmt.Errorf("gagal memulai transaksi database: %v", err)
    }
    defer tx.Rollback()

    // 1. Kalkulasi Total dan Detail
    var totalTrans float64
    var details []models.BeliDetail

    for _, d := range req.Details {
        // Validasi Barang (Optional, cek apakah barang exists)
        // _, err := s.barangRepo.GetByID(d.BarangID)
        // if err != nil { return nil, fmt.Errorf("barang tidak ditemukan") }

        subtotal := float64(d.Qty) * d.Harga
        totalTrans += subtotal
        
        details = append(details, models.BeliDetail{
            BarangID: d.BarangID,
            Qty:      d.Qty,
            Harga:    d.Harga,
            Subtotal: subtotal,
        })
    }

    // 2. Auto Generate No Faktur
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

    // 3. Simpan Transaksi ke Database
    if err := s.repo.Create(tx, header, details); err != nil {
        return nil, fmt.Errorf("gagal membuat transaksi: %v", err)
    }

    // 4. Update Stok & Buat History
    for _, d := range details {
        // Ambil stok sebelum update
        currentStok, err := s.stokRepo.GetByBarangID(d.BarangID)
        var stokSebelum int
        if err != nil || currentStok == nil {
            stokSebelum = 0
        } else {
            stokSebelum = currentStok.StokAkhir
        }
        
        // Tambah Stok
        if err := s.stokRepo.CreateOrUpdate(tx, d.BarangID, d.Qty); err != nil {
             return nil, fmt.Errorf("gagal memperbarui stok barang ID %d: %v", d.BarangID, err)
        }
        
        // Hitung stok sesudah
        stokSesudah := stokSebelum + d.Qty
        
        // Catat History
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

    // 5. Commit Transaksi
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("gagal commit transaksi: %v", err)
    }

    return header, nil
}
