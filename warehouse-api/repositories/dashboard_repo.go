package repositories

import (
    "database/sql"
    "warehouse-api/models"
)

type DashboardRepository interface {
    GetStats() (*models.DashboardStats, error)
}

type dashboardRepository struct {
    db *sql.DB
}

func NewDashboardRepository(db *sql.DB) DashboardRepository {
    return &dashboardRepository{db}
}

func (r *dashboardRepository) GetStats() (*models.DashboardStats, error) {
    stats := &models.DashboardStats{}

    // 1. Total User
    err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUser)
    if err != nil { return nil, err }

    // 2. Total Barang
    err = r.db.QueryRow("SELECT COUNT(*) FROM master_barang").Scan(&stats.TotalBarang)
    if err != nil { return nil, err }

    // 3. Total Stok
    // Asumsi tabel mstok memiliki semua barang
    err = r.db.QueryRow("SELECT COALESCE(SUM(stok_akhir), 0) FROM mstok").Scan(&stats.TotalStok)
    if err != nil { return nil, err }

    // 4. Total Nilai Aset (Harga Beli * Stok)
    // Join master_barang dan mstok
    queryAset := `
        SELECT COALESCE(SUM(b.harga_beli * s.stok_akhir), 0)
        FROM master_barang b
        JOIN mstok s ON b.id = s.barang_id
    `
    err = r.db.QueryRow(queryAset).Scan(&stats.TotalNilaiAset)
    if err != nil { return nil, err }

    // 5. Top 5 Barang Terlaris (Berdasarkan table jual_detail)
    queryTop := `
        SELECT b.nama_barang, COALESCE(SUM(d.qty), 0) as total_terjual
        FROM jual_detail d
        JOIN master_barang b ON d.barang_id = b.id
        GROUP BY b.id, b.nama_barang
        ORDER BY total_terjual DESC
        LIMIT 5
    `
    rows, err := r.db.Query(queryTop)
    if err != nil {
        // Log error di sini jika perlu, tapi kembalikan slice kosong agar dashboard tetap jalan
        stats.TopSellingProducts = []models.TopProduct{}
    } else {
        defer rows.Close()
        for rows.Next() {
            var p models.TopProduct
            if err := rows.Scan(&p.NamaBarang, &p.TotalTerjual); err != nil {
                continue
            }
            stats.TopSellingProducts = append(stats.TopSellingProducts, p)
        }
    }
    
    return stats, nil
}
