package models

type Barang struct {
	ID         int     `json:"id"`
	KodeBarang string  `json:"kode_barang"`
	NamaBarang string  `json:"nama_barang"`
	Deskripsi  string  `json:"deskripsi"`
	Satuan     string  `json:"satuan"`
	HargaBeli  float64 `json:"harga_beli"`
	HargaJual  float64 `json:"harga_jual"`
}

type BarangWithStok struct {
	Barang
	Stok int `json:"stok"`
}

type CreateBarangRequest struct {
	NamaBarang string  `json:"nama_barang"`
	Deskripsi  string  `json:"deskripsi"`
	Satuan     string  `json:"satuan"`
	HargaBeli  float64 `json:"harga_beli"`
	HargaJual  float64 `json:"harga_jual"`
}
