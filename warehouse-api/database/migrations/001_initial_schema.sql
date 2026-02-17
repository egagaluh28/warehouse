-- Table User
CREATE TABLE IF NOT EXISTS users (
 id SERIAL PRIMARY KEY,
 username VARCHAR(100) UNIQUE NOT NULL,
 password VARCHAR(255) NOT NULL,
 email VARCHAR(150) UNIQUE NOT NULL,
 full_name VARCHAR(200) NOT NULL,
 role VARCHAR(50) DEFAULT 'staff',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Master Barang
CREATE TABLE IF NOT EXISTS master_barang (
 id SERIAL PRIMARY KEY,
 kode_barang VARCHAR(50) UNIQUE NOT NULL,
 nama_barang VARCHAR(200) NOT NULL,
 deskripsi TEXT,
 satuan VARCHAR(50) NOT NULL, 
 harga_beli DECIMAL(15,2) DEFAULT 0,
 harga_jual DECIMAL(15,2) DEFAULT 0,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Stok
CREATE TABLE IF NOT EXISTS mstok (
 id SERIAL PRIMARY KEY,
 barang_id INTEGER REFERENCES master_barang(id),
 stok_akhir INTEGER DEFAULT 0,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table History Stok
CREATE TABLE IF NOT EXISTS history_stok (
 id SERIAL PRIMARY KEY,
 barang_id INTEGER REFERENCES master_barang(id),
 user_id INTEGER REFERENCES users(id),
 jenis_transaksi VARCHAR(50) NOT NULL, -- 'masuk', 'keluar', 'adjustment'
 jumlah INTEGER NOT NULL,
 stok_sebelum INTEGER NOT NULL,
 stok_sesudah INTEGER NOT NULL,
 keterangan TEXT,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Pembelian Header
CREATE TABLE IF NOT EXISTS beli_header (
 id SERIAL PRIMARY KEY,
 no_faktur VARCHAR(100) UNIQUE NOT NULL,
 supplier VARCHAR(200) NOT NULL,
 total DECIMAL(15,2) DEFAULT 0,
 user_id INTEGER REFERENCES users(id),
 status VARCHAR(50) DEFAULT 'selesai',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Pembelian Detail
CREATE TABLE IF NOT EXISTS beli_detail (
 id SERIAL PRIMARY KEY,
 beli_header_id INTEGER REFERENCES beli_header(id),
 barang_id INTEGER REFERENCES master_barang(id),
 qty INTEGER NOT NULL,
 harga DECIMAL(15,2) NOT NULL,
 subtotal DECIMAL(15,2) NOT NULL
);

-- Table Penjualan Header
CREATE TABLE IF NOT EXISTS jual_header (
 id SERIAL PRIMARY KEY,
 no_faktur VARCHAR(100) UNIQUE NOT NULL,
 customer VARCHAR(200) NOT NULL,
 total DECIMAL(15,2) DEFAULT 0,
 user_id INTEGER REFERENCES users(id),
 status VARCHAR(50) DEFAULT 'selesai',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Penjualan Detail
CREATE TABLE IF NOT EXISTS jual_detail (
 id SERIAL PRIMARY KEY,
 jual_header_id INTEGER REFERENCES jual_header(id),
 barang_id INTEGER REFERENCES master_barang(id),
 qty INTEGER NOT NULL,
 harga DECIMAL(15,2) NOT NULL,
 subtotal DECIMAL(15,2) NOT NULL
);
