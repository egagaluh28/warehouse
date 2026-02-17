package models

type DashboardStats struct {
    TotalUser          int     `json:"total_user"`
    TotalBarang        int     `json:"total_barang"`
    TotalStok          int     `json:"total_stok"`
    TotalNilaiAset     float64 `json:"total_nilai_aset"`
    TopSellingProducts []TopProduct `json:"top_selling_products"`
}

type TopProduct struct {
    NamaBarang   string `json:"nama_barang"`
    TotalTerjual int    `json:"total_terjual"`
}
