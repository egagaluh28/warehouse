package models

import "time"

type Stok struct {
	ID        int       `json:"id"`
	BarangID  int       `json:"barang_id"`
	StokAkhir int       `json:"stok_akhir"`
	UpdatedAt time.Time `json:"updated_at"`
	Barang    *Barang   `json:"barang,omitempty"`
}


