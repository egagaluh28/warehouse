# 📦 Warehouse Inventory API

> REST API untuk sistem manajemen inventaris gudang yang komprehensif, dibangun dengan Go dan PostgreSQL.

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Test Coverage](https://img.shields.io/badge/coverage-75%25-brightgreen?style=flat)](TEST_COVERAGE_REPORT.md)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ⚡ Quick Start

```bash
# Clone repository
git clone https://github.com/username/warehouse-api.git
cd warehouse-api

# Install dependencies
go mod download

# Setup database
psql -U postgres -c "CREATE DATABASE warehouse;"
psql -U postgres -d warehouse -f database/migrations/001_initial_schema.sql
psql -U postgres -d warehouse -f database/migrations/002_fix_admin_password.sql

# Create .env file (sesuaikan dengan konfigurasi Anda)
echo "DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=warehouse
JWT_SECRET=your-secret-key" > .env

# Run application
go run main.go

# API ready at: http://localhost:8080
# Swagger docs: http://localhost:8080/swagger/index.html
# Default login: admin / admin123
```

## 📋 Daftar Isi

- [Quick Start](#-quick-start)
- [Fitur](#-fitur)
- [Teknologi](#-teknologi)
- [Prasyarat](#-prasyarat)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [Dokumentasi API](#-dokumentasi-api)
- [Struktur Project](#-struktur-project)
- [Database Schema](#-database-schema)
- [Response Format](#-response-format)
- [Security Features](#-security-features)
- [Testing](#-testing)
- [Troubleshooting](#-troubleshooting)
- [Performance Tips](#-performance-tips)
- [Development](#-development)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)
- [Contact](#-contact)

## ✨ Fitur

- 🔐 **Autentikasi & Autorisasi** - JWT-based authentication dengan role management (Admin/Staff)
- 📊 **Dashboard Analytics** - Statistik real-time untuk monitoring gudang
- 📦 **Manajemen Barang** - CRUD lengkap dengan pencarian dan pagination
- 📈 **Tracking Stok** - Monitoring stok real-time dengan history lengkap
- 🛒 **Transaksi Pembelian** - Pencatatan pembelian dengan multi-item support
- 💰 **Transaksi Penjualan** - Pencatatan penjualan dengan validasi stok otomatis
- 📝 **Audit Trail** - History lengkap semua pergerakan stok
- 🚀 **Rate Limiting** - Proteksi API dari abuse
- 📖 **Swagger Documentation** - API documentation interaktif

## 🛠 Teknologi

- **Backend Framework:** Go (Golang) 1.24
- **Database:** PostgreSQL 14+
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **API Documentation:** Swagger/OpenAPI 2.0
- **Architecture Pattern:** Clean Architecture (Handler → Service → Repository)

## 📦 Prasyarat

Pastikan sudah terinstall:

- [Go](https://golang.org/dl/) 1.22 atau lebih tinggi
- [PostgreSQL](https://www.postgresql.org/download/) 14 atau lebih tinggi
- [Git](https://git-scm.com/downloads)

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/username/warehouse-api.git
cd warehouse-api
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE warehouse;
```

Jalankan migration untuk membuat skema tabel:

```bash
psql -U postgres -d warehouse -f database/migrations/001_initial_schema.sql
psql -U postgres -d warehouse -f database/migrations/002_fix_admin_password.sql
```

### 4. Seed Data (Opsional)

Untuk data testing:

```bash
go run cmd/seeder/main.go
```

## ⚙ Konfigurasi

Buat file `.env` di root project:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=warehouse

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Configuration (optional)
PORT=8080
```

## 🏃 Menjalankan Aplikasi

### Development Mode

```bash
go run main.go
```

### Build & Run

```bash
go build -o warehouse-api
./warehouse-api
```

Server akan berjalan di: **http://localhost:8080**

### Akses Dokumentasi Swagger

Buka browser dan akses: **http://localhost:8080/swagger/index.html**

## 📖 Dokumentasi API

### Base URL

```
http://localhost:8080/api
```

### Authentication

Semua endpoint (kecuali `/login`) memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

### Default Admin Credentials

```json
{
  "username": "admin",
  "password": "admin123"
}
```

---

## 🔐 Authentication

### Login

**POST** `/api/login`

Mendapatkan JWT token untuk autentikasi.

**Request Body:**

```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "full_name": "Administrator",
      "role": "admin"
    }
  }
}
```

### Register User

**POST** `/api/register`

Mendaftarkan user baru (hanya untuk admin).

**Headers:**

```
Authorization: Bearer <admin-token>
```

**Request Body:**

```json
{
  "username": "staff01",
  "password": "password123",
  "full_name": "Staff Member",
  "email": "staff@example.com",
  "role": "staff"
}
```

**Response:**

```json
{
  "success": true,
  "message": "User berhasil didaftarkan",
  "data": {
    "id": 2,
    "username": "staff01",
    "full_name": "Staff Member"
  }
}
```

---

## 📊 Dashboard

### Get Dashboard Statistics

**GET** `/api/dashboard`

Mengambil ringkasan statistik gudang.

**Response:**

```json
{
  "success": true,
  "data": {
    "total_barang": 150,
    "total_stok": 5420,
    "nilai_aset": 125500000,
    "top_selling": [
      {
        "barang_id": 1,
        "nama_barang": "Laptop ASUS ROG",
        "total_terjual": 45
      }
    ]
  }
}
```

---

## 📦 Manajemen Barang

### Get All Barang

**GET** `/api/barang`

Mengambil daftar barang dengan pencarian dan pagination.

**Query Parameters:**

- `search` (string): Pencarian berdasarkan nama/kode
- `page` (integer): Nomor halaman (default: 1)
- `limit` (integer): Jumlah per halaman (default: 10)
- `sort_by` (string): Field untuk sorting (nama, harga_beli, harga_jual)
- `order` (string): asc/desc (default: asc)

**Example:**

```
GET /api/barang?search=laptop&page=1&limit=10&sort_by=harga_jual&order=desc
```

**Response:**

```json
{
  "success": true,
  "message": "Data barang berhasil diambil",
  "data": [
    {
      "id": 1,
      "kode_barang": "BRG-001",
      "nama_barang": "Laptop ASUS ROG",
      "satuan": "Unit",
      "harga_beli": 12000000,
      "harga_jual": 15000000,
      "deskripsi": "Laptop gaming high-end"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 150
  }
}
```

### Get Barang by ID

**GET** `/api/barang/{id}`

**Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "kode_barang": "BRG-001",
    "nama_barang": "Laptop ASUS ROG",
    "satuan": "Unit",
    "harga_beli": 12000000,
    "harga_jual": 15000000,
    "deskripsi": "Laptop gaming high-end"
  }
}
```

### Create Barang

**POST** `/api/barang`

`kode_barang` akan dibuat otomatis oleh sistem.

**Request Body:**

```json
{
  "nama_barang": "Mouse Logitech G502",
  "satuan": "Pcs",
  "harga_beli": 500000,
  "harga_jual": 750000,
  "deskripsi": "Gaming mouse dengan RGB"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Barang berhasil ditambahkan",
  "data": {
    "id": 2,
    "kode_barang": "BRG-002",
    "nama_barang": "Mouse Logitech G502"
  }
}
```

### Update Barang

**PUT** `/api/barang/{id}`

**Request Body:**

```json
{
  "nama_barang": "Mouse Logitech G502 Hero",
  "harga_jual": 800000
}
```

### Delete Barang

**DELETE** `/api/barang/{id}`

**Response:**

```json
{
  "success": true,
  "message": "Barang berhasil dihapus"
}
```

---

## 📊 Manajemen Stok

### Get All Stok

**GET** `/api/stok`

Mengambil data stok semua barang.

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "barang_id": 1,
      "nama_barang": "Laptop ASUS ROG",
      "qty": 25,
      "last_update": "2026-02-17T14:30:00Z"
    }
  ]
}
```

### Get Stok by Barang ID

**GET** `/api/stok/{id}`

**Response:**

```json
{
  "success": true,
  "data": {
    "barang_id": 1,
    "nama_barang": "Laptop ASUS ROG",
    "qty": 25,
    "last_update": "2026-02-17T14:30:00Z"
  }
}
```

### Get History Stok

**GET** `/api/history-stok`

**GET** `/api/history-stok/{id}` (filter by barang_id)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "barang_id": 1,
      "nama_barang": "Laptop ASUS ROG",
      "tipe": "masuk",
      "qty": 10,
      "keterangan": "Pembelian - FK-001",
      "created_at": "2026-02-17T10:00:00Z"
    }
  ]
}
```

---

## 🛒 Transaksi Pembelian

### Create Pembelian

**POST** `/api/pembelian`

Mencatat transaksi pembelian (stok masuk).

**Request Body:**

```json
{
  "no_faktur": "FK-001",
  "supplier": "PT Tech Supplier",
  "user_id": 1,
  "details": [
    {
      "barang_id": 1,
      "qty": 10,
      "harga": 12000000
    },
    {
      "barang_id": 2,
      "qty": 5,
      "harga": 500000
    }
  ]
}
```

**Response:**

```json
{
  "success": true,
  "message": "Pembelian berhasil dicatat",
  "data": {
    "id": 1,
    "no_faktur": "FK-001",
    "total": 122500000
  }
}
```

### Get All Pembelian

**GET** `/api/pembelian`

**Query Parameters:**

- `start_date` (string): Filter dari tanggal (YYYY-MM-DD)
- `end_date` (string): Filter sampai tanggal (YYYY-MM-DD)

**Example:**

```
GET /api/pembelian?start_date=2026-02-01&end_date=2026-02-28
```

### Get Pembelian by ID

**GET** `/api/pembelian/{id}`

**Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "no_faktur": "FK-001",
    "supplier": "PT Tech Supplier",
    "tanggal": "2026-02-17T10:00:00Z",
    "total": 122500000,
    "details": [
      {
        "barang_id": 1,
        "nama_barang": "Laptop ASUS ROG",
        "qty": 10,
        "harga": 12000000,
        "subtotal": 120000000
      }
    ]
  }
}
```

---

## 💰 Transaksi Penjualan

### Create Penjualan

**POST** `/api/penjualan`

Mencatat transaksi penjualan (stok keluar).

**Request Body:**

```json
{
  "no_faktur": "INV-001",
  "customer": "PT Client ABC",
  "user_id": 1,
  "details": [
    {
      "barang_id": 1,
      "qty": 2,
      "harga": 15000000
    }
  ]
}
```

**Response:**

```json
{
  "success": true,
  "message": "Penjualan berhasil dicatat",
  "data": {
    "id": 1,
    "no_faktur": "INV-001",
    "total": 30000000
  }
}
```

### Get All Penjualan

**GET** `/api/penjualan`

**Query Parameters:**

- `start_date` (string): Filter dari tanggal (YYYY-MM-DD)
- `end_date` (string): Filter sampai tanggal (YYYY-MM-DD)

### Get Penjualan by ID

**GET** `/api/penjualan/{id}`

---

## 🏗 Struktur Project

```
warehouse-api/
├── cmd/
│   └── seeder/
│       └── main.go              # Data seeder
├── config/
│   └── database.go              # Database connection
├── database/
│   ├── migrations/              # SQL migration files
│   │   ├── 001_initial_schema.sql
│   │   └── 002_fix_admin_password.sql
│   └── seeders/                 # Go seeders
│       ├── barang_seeder.go
│       └── user_seeder.go
├── docs/                        # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── handlers/                    # HTTP handlers / controllers
│   ├── barang_handler.go
│   ├── dashboard_handler.go
│   ├── pembelian_handler.go
│   ├── penjualan_handler.go
│   ├── stok_handler.go
│   └── user_handler.go
├── middleware/                  # HTTP middlewares
│   ├── auth.go                  # JWT authentication
│   ├── logger.go                # Request logging
│   └── ratelimit.go             # Rate limiting
├── models/                      # Data models / entities
│   ├── barang.go
│   ├── dashboard.go
│   ├── history_stok.go
│   ├── pembelian.go
│   ├── penjualan.go
│   ├── response.go
│   ├── stok.go
│   └── user.go
├── repositories/                # Data access layer
│   ├── barang_repo.go
│   ├── dashboard_repo.go
│   ├── pembelian_repo.go
│   ├── penjualan_repo.go
│   ├── stok_repo.go
│   └── user_repo.go
├── services/                    # Business logic layer
│   ├── pembelian_service.go
│   └── penjualan_service.go
├── utils/                       # Utility functions
│   ├── generator.go
│   └── response_helper.go
├── .env                         # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── main.go                      # Application entry point
└── README.md
```

## 💾 Database Schema

### Tabel Users

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Tabel Barang

```sql
CREATE TABLE barang (
    id SERIAL PRIMARY KEY,
    kode_barang VARCHAR(50) UNIQUE NOT NULL,
    nama_barang VARCHAR(100) NOT NULL,
    satuan VARCHAR(20) NOT NULL,
    harga_beli DECIMAL(15,2) NOT NULL,
    harga_jual DECIMAL(15,2) NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Tabel Stok

```sql
CREATE TABLE stok (
    barang_id INTEGER PRIMARY KEY REFERENCES barang(id),
    qty INTEGER NOT NULL DEFAULT 0,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Tabel History Stok

```sql
CREATE TABLE history_stok (
    id SERIAL PRIMARY KEY,
    barang_id INTEGER REFERENCES barang(id),
    tipe VARCHAR(10) NOT NULL, -- 'masuk' or 'keluar'
    qty INTEGER NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Tabel Pembelian

```sql
CREATE TABLE pembelian (
    id SERIAL PRIMARY KEY,
    no_faktur VARCHAR(50) UNIQUE NOT NULL,
    tanggal TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    supplier VARCHAR(100),
    total DECIMAL(15,2),
    user_id INTEGER REFERENCES users(id)
);
```

### Tabel Pembelian Detail

```sql
CREATE TABLE pembelian_detail (
    id SERIAL PRIMARY KEY,
    pembelian_id INTEGER REFERENCES pembelian(id) ON DELETE CASCADE,
    barang_id INTEGER REFERENCES barang(id),
    qty INTEGER NOT NULL,
    harga DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL
);
```

### Tabel Penjualan & Penjualan Detail

Struktur mirip dengan tabel pembelian.

## 📝 Response Format

Semua response API mengikuti format standar:

### Success Response

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": {
    /* data object atau array */
  },
  "meta": {
    /* pagination info jika ada */
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Pesan error",
  "data": null
}
```

## 🔒 Security Features

- ✅ JWT-based authentication
- ✅ Password hashing dengan bcrypt
- ✅ Role-based access control (Admin/Staff)
- ✅ Rate limiting untuk mencegah abuse
- ✅ CORS configuration
- ✅ SQL injection protection (menggunakan prepared statements)
- ✅ Input validation

## 🧪 Testing

### Test Coverage

Project ini memiliki **75% code coverage** dengan **40+ test cases** yang mencakup unit tests dan integration tests.

### Test Structure

```
test/
├── unit/                          # Unit Tests (9 files)
│   ├── generator_test.go         # Utils code generation tests
│   ├── response_helper_test.go   # HTTP response helper tests
│   ├── user_service_test.go      # User authentication tests
│   ├── barang_handler_test.go    # Barang handler tests
│   ├── middleware_test.go        # Auth, logging, rate limit tests
│   └── ...                       # Dan lainnya
│
└── integration/                   # Integration Tests (2 files)
    ├── api_test.go               # End-to-end API tests
    └── database_test.go          # Database CRUD tests
```

### Running Tests

**Run All Tests:**

```bash
# Semua tests (unit + integration)
go test ./test/...

# Dengan verbose output
go test -v ./test/...
```

**Run Unit Tests Only:**

```bash
# Hanya unit tests
go test -v ./test/unit

# Dengan coverage report
go test -v -cover ./test/unit
```

**Run Integration Tests Only:**

```bash
# Hanya integration tests (memerlukan database)
go test -v ./test/integration
```

**Generate Coverage Report:**

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./test/...

# View coverage sebagai HTML
go tool cover -html=coverage.out -o coverage.html
```

### Test Prerequisites

**Dependencies:**

```bash
# Install testing dependencies
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
```

**Setup Test Database (untuk integration tests):**

```sql
CREATE DATABASE warehouse_test;
```

### Postman Collection

Untuk manual testing API, gunakan Postman Collection yang tersedia di:

```
warehouse_api_postman_collection.json
```

Import file tersebut ke Postman untuk testing API endpoints.

### Dokumentasi Testing Lengkap

Untuk panduan testing yang lebih detail, lihat:

- **[TESTING.md](TESTING.md)** - Comprehensive testing guide
- **[TEST_COVERAGE_REPORT.md](TEST_COVERAGE_REPORT.md)** - Detailed coverage metrics

---

## 🐛 Troubleshooting

### Database Connection Issues

**Error: "connection refused"**

```bash
# Pastikan PostgreSQL berjalan
# Windows:
Get-Service -Name postgresql*

# Linux/Mac:
sudo systemctl status postgresql
```

**Error: "database does not exist"**

```bash
# Buat database terlebih dahulu
psql -U postgres -c "CREATE DATABASE warehouse;"
```

### Migration Errors

**Error: "relation already exists"**

```bash
# Drop database dan buat ulang
psql -U postgres -c "DROP DATABASE warehouse;"
psql -U postgres -c "CREATE DATABASE warehouse;"

# Jalankan migrations lagi
psql -U postgres -d warehouse -f database/migrations/001_initial_schema.sql
```

### JWT Token Issues

**Error: "invalid token"**

- Pastikan JWT_SECRET di `.env` sama dengan yang digunakan saat generate token
- Token mungkin sudah expired (default 24 jam)
- Login ulang untuk mendapat token baru

### Common Test Failures

**Integration tests gagal:**

```bash
# Pastikan test database sudah dibuat
psql -U postgres -c "CREATE DATABASE warehouse_test;"

# Pastikan .env atau environment variables sudah diset
```

---

## 🚀 Performance Tips

### Database Optimization

```sql
-- Tambahkan index untuk pencarian yang sering digunakan
CREATE INDEX idx_barang_nama ON master_barang(nama_barang);
CREATE INDEX idx_barang_kode ON master_barang(kode_barang);
CREATE INDEX idx_history_stok_barang ON history_stok(barang_id);
```

### Pagination Best Practices

Selalu gunakan pagination untuk query yang return banyak data:

```
GET /api/barang?page=1&limit=50
```

### Rate Limiting

API sudah dilengkapi rate limiting. Default:

- **100 requests per menit** per IP address

Untuk production, sesuaikan di `middleware/ratelimit.go`.

---

## 🔧 Development

### Hot Reload (Recommended)

Install **Air** untuk hot reload saat development:

```bash
# Install Air
go install github.com/air-verse/air@latest

# Run dengan hot reload
air
```

### Code Quality Tools

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (optional - install golangci-lint)
golangci-lint run
```

### Generate Swagger Docs

Jika ada perubahan pada API documentation:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swag init

# Swagger akan available di /swagger/index.html
```

---

## 📦 Deployment

### Build for Production

```bash
# Build binary
go build -o warehouse-api

# Build dengan optimizations
go build -ldflags="-s -w" -o warehouse-api

# Cross-compile untuk Linux
GOOS=linux GOARCH=amd64 go build -o warehouse-api-linux
```

### Environment Variables Production

Pastikan set environment variables untuk production:

```env
DB_HOST=your-production-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=warehouse

JWT_SECRET=your-super-secure-random-secret-key-min-32-characters

PORT=8080
```

### Running in Production

```bash
# Direct run
./warehouse-api

# With systemd (Linux)
sudo systemctl start warehouse-api

# With PM2 (if using Node.js ecosystem)
pm2 start ./warehouse-api --name warehouse-api

# With Docker (create Dockerfile first)
docker build -t warehouse-api .
docker run -p 8080:8080 warehouse-api
```

---

## 🤝 Contributing

Contributions are welcome! Untuk berkontribusi:

1. Fork repository ini
2. Buat branch baru (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

### Coding Standards

- Ikuti [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Tulis unit tests untuk fitur baru
- Pastikan semua tests passing sebelum PR
- Update dokumentasi jika diperlukan

---

## 📄 License

[MIT License](LICENSE)

## 👥 Contributors

- **Your Name** - Initial work

## 📧 Contact

Untuk pertanyaan atau saran, silakan hubungi:

- Email: your.email@example.com
- GitHub: [@yourusername](https://github.com/yourusername)

---

**Made with ❤️ using Go**
