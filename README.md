# Warehouse Management System (Backend + Frontend)

Monorepo ini berisi:

- **Backend API (Go + PostgreSQL)**: folder `warehouse-api/`
- **Frontend Web (Next.js + NextUI + Tailwind)**: folder `warehouse-fe/`

## Fitur yang sudah ada

- **Auth & Role**: login JWT + proteksi route (Admin/Staff)
- **Dashboard**: ringkasan statistik gudang
- **Barang**:
  - CRUD barang + search + pagination
  - `kode_barang` **dibuat otomatis** oleh backend saat create
  - Endpoint tambahan **barang + stok**: `GET /api/barang/stok`
- **Stok**: monitoring stok + history stok
- **Transaksi**:
  - Pembelian (stok masuk) multi-item
  - Penjualan (stok keluar) dengan validasi stok
- **Invoice PDF (Frontend)**:
  - Download invoice PDF untuk **Pembelian** dan **Penjualan** dari tabel (action baru)
- **Dokumentasi API**: Swagger UI
- **Keamanan dasar**: bcrypt, JWT, rate limit, CORS

## Quick Start (Local)

### 1) Backend

Lihat panduan lengkap di [warehouse-api/README.md](warehouse-api/README.md).

Ringkasnya:

```bash
cd warehouse-api

# buat DB + migrate
psql -U postgres -c "CREATE DATABASE warehouse;"
psql -U postgres -d warehouse -f database/migrations/001_initial_schema.sql
psql -U postgres -d warehouse -f database/migrations/002_fix_admin_password.sql

# konfigurasi env
copy .env.example .env

# run
go run main.go
```

Backend URL:

- API: `http://localhost:8080/api`
- Swagger: `http://localhost:8080/swagger/index.html`

Default login:

- `admin / admin`

### 2) Frontend

Lihat panduan lengkap di [warehouse-fe/README.md](warehouse-fe/README.md).

Ringkasnya:

```bash
cd warehouse-fe
npm install

# set API base URL
copy .env.local .env.local

npm run dev
```

Frontend URL:

- Web: `http://localhost:3000`

## Struktur Folder

```text
.
├─ warehouse-api/   # Go REST API + Swagger + tests
└─ warehouse-fe/    # Next.js App Router (web)
```
