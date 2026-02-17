package models

import "time"

type JualHeader struct {
	ID        int          `json:"id"`
	NoFaktur  string       `json:"no_faktur"`
	Customer  string       `json:"customer"`
	Total     float64      `json:"total"`
	UserID    int          `json:"user_id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	User      *User        `json:"user,omitempty"`
	Details   []JualDetail `json:"details,omitempty"`
}

type JualDetail struct {
	ID           int     `json:"id"`
	JualHeaderID int     `json:"jual_header_id"`
	BarangID     int     `json:"barang_id"`
	Qty          int     `json:"qty"`
	Harga        float64 `json:"harga"` // Harga Jual
	Subtotal     float64 `json:"subtotal"`
	Barang       *Barang `json:"barang,omitempty"`
}

type CreatePenjualanRequest struct {
	NoFaktur string                  `json:"no_faktur"` // Optional, or generated
	Customer string                  `json:"customer"`
	UserID   int                     `json:"user_id"`
	Details  []CreatePenjualanDetail `json:"details"`
}

type CreatePenjualanDetail struct {
	BarangID int     `json:"barang_id"`
	Qty      int     `json:"qty"`
	Harga    float64 `json:"harga"`
}
