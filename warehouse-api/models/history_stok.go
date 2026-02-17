package models

import "time"

type HistoryStok struct {
	ID            int       `json:"id"`
	BarangID      int       `json:"barang_id"`
	UserID        int       `json:"user_id"`
	JenisTransaksi string   `json:"jenis_transaksi"`
	Jumlah        int       `json:"jumlah"`
	StokSebelum   int       `json:"stok_sebelum"`
	StokSesudah   int       `json:"stok_sesudah"`
	Keterangan    string    `json:"keterangan"`
	CreatedAt     time.Time `json:"created_at"`
	Barang        *Barang   `json:"barang,omitempty"`
	User          *User     `json:"user,omitempty"`
}
