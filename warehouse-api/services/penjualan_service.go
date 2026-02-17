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
    tx, err := s.db.Begin()
    if err != nil {
        return nil, fmt.Errorf("gagal memulai transaksi database: %v", err)
    }
    defer tx.Rollback()

    var totalTrans float64
    var details []models.JualDetail

    // 1. Validasi & Kalkulasi Detail Barang
    for _, d := range req.Details {
        currentStok, err := s.stokRepo.GetByBarangID(d.BarangID)
        if err != nil {
            return nil, fmt.Errorf("barang ID %d tidak ditemukan dalam stok", d.BarangID)
        }
        
        if currentStok.StokAkhir < d.Qty {
            return nil, fmt.Errorf("stok tidak mencukupi untuk Barang ID %d. Tersedia: %d, Diminta: %d", d.BarangID, currentStok.StokAkhir, d.Qty)
        }

        subtotal := float64(d.Qty) * d.Harga
        totalTrans += subtotal
        
        details = append(details, models.JualDetail{
            BarangID: d.BarangID,
            Qty:      d.Qty,
            Harga:    d.Harga,
            Subtotal: subtotal,
        })
    }

    // 2. Auto Generate No Faktur
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

    // 3. Simpan Header & Detail Transaksi
    if err := s.repo.Create(tx, header, details); err != nil {
        return nil, fmt.Errorf("gagal membuat transaksi: %v", err)
    }

    // 4. Update Stok & Catat History
    for _, d := range details {
        // Kurangi Stok
        if err := s.stokRepo.CreateOrUpdate(tx, d.BarangID, -d.Qty); err != nil {
             return nil, fmt.Errorf("gagal memperbarui stok barang ID %d: %v", d.BarangID, err)
        }
        
        // Catat History
        
        history := &models.HistoryStok{
            BarangID:       d.BarangID,
            UserID:         req.UserID,
            JenisTransaksi: "keluar",
            Jumlah:         d.Qty,
            StokSebelum:    0, 
            StokSesudah:    0, 
            Keterangan:     "Penjualan " + header.NoFaktur,
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
