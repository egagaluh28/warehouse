# Warehouse Inventory API (Backend)

REST API untuk sistem manajemen inventaris gudang, dibangun dengan Go + PostgreSQL.

## Fitur yang sudah ada

- JWT Authentication + role (Admin/Staff)
- Dashboard statistik
- Barang (CRUD) + search + pagination
  - `kode_barang` **dibuat otomatis** oleh sistem saat create
  - Endpoint tambahan barang + stok: `GET /api/barang/stok`
- Stok + history stok
- Transaksi Pembelian (stok masuk) multi-item
- Transaksi Penjualan (stok keluar) dengan validasi stok
- Swagger UI
- Middleware: logger, CORS, rate limiting

## Prasyarat

- Go 1.22+ (project sudah berjalan pada Go 1.24)
- PostgreSQL 14+

## Konfigurasi

Copy file env:

```bash
cd warehouse-api
copy .env.example .env
```

Isi minimal:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=warehouse
JWT_SECRET=your-secret-key
PORT=8080
```

## Setup Database

```bash
cd warehouse-api

psql -U postgres -c "CREATE DATABASE warehouse;"
psql -U postgres -d warehouse -f database/migrations/001_initial_schema.sql
psql -U postgres -d warehouse -f database/migrations/002_fix_admin_password.sql

# optional seed
go run cmd/seeder/main.go
```

## Menjalankan Backend

```bash
cd warehouse-api
go run main.go
```

- API Base URL: `http://localhost:8080/api`
- Swagger: `http://localhost:8080/swagger/index.html`

Default login:

- `admin / admin`

## Endpoint Ringkas

Base path: `/api`

- Auth: `POST /login`, `POST /register` (register untuk admin)
- Dashboard: `GET /dashboard`
- Barang:
  - `GET /barang` (list)
  - `GET /barang/{id}`
  - `POST /barang` (tanpa `kode_barang`, dibuat otomatis)
  - `PUT /barang/{id}`
  - `DELETE /barang/{id}`
  - `GET /barang/stok` (list barang + stok)
- Stok: `GET /stok`, `GET /stok/{id}`
- History stok: `GET /history-stok`, `GET /history-stok/{id}` (filter by barang_id)
- Pembelian: `GET /pembelian`, `GET /pembelian/{id}`, `POST /pembelian`
- Penjualan: `GET /penjualan`, `GET /penjualan/{id}`, `POST /penjualan`

## Testing

Jalankan semua test:

```bash
cd warehouse-api
go test -v ./...
```

Catatan:

- Test integration butuh koneksi database (lihat `test/integration`).
- Beberapa test unit memang `SKIP` (misalnya yang butuh mocking transaksi DB).

### Hasil test terakhir (local)

Command:

```bash
go test -v ./...
```

Output ringkas:

```text
ok      warehouse-api/test/integration  1.106s
ok      warehouse-api/test/unit         1.179s
```
