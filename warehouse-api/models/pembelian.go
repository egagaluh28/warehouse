package models

import "time"

type BeliHeader struct {
	ID        int          `json:"id"`
	NoFaktur  string       `json:"no_faktur"`
	Supplier  string       `json:"supplier"`
	Total     float64      `json:"total"`
	UserID    int          `json:"user_id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	User      *User        `json:"user,omitempty"`
	Details   []BeliDetail `json:"details,omitempty"`
}

type BeliDetail struct {
	ID           int     `json:"id"`
	BeliHeaderID int     `json:"beli_header_id"`
	BarangID     int     `json:"barang_id"`
	Qty          int     `json:"qty"`
	Harga        float64 `json:"harga"`
	Subtotal     float64 `json:"subtotal"`
	Barang       *Barang `json:"barang,omitempty"`
}

type CreatePembelianRequest struct {
	NoFaktur string                  `json:"no_faktur"` // Optional, or generated
	Supplier string                  `json:"supplier"`
	UserID   int                     `json:"user_id"`
	Details  []CreatePembelianDetail `json:"details"`
}

type CreatePembelianDetail struct {
	BarangID int     `json:"barang_id"`
	Qty      int     `json:"qty"`
	Harga    float64 `json:"harga"`
}
