# Warehouse Web (Frontend)

Frontend web untuk sistem Warehouse Management.

Tech stack utama:

- Next.js (App Router)
- React
- NextUI / HeroUI
- TailwindCSS
- Axios (API client) + JWT via `localStorage`
- PDF invoice: `jspdf` + `jspdf-autotable`

## Fitur yang sudah ada

- Login & protected routes
- Dashboard
- Barang (list, create, update, delete) + detail (fetch by `getById`)
- Stok + history stok
- Pembelian + Penjualan (list + detail)
- **Download Invoice PDF** untuk Pembelian & Penjualan (action baru di tabel)
- UI tabel terstandarisasi via komponen `DataTable`

## Prasyarat

- Node.js 18+ (disarankan 20+)
- Backend berjalan (lihat [../warehouse-api/README.md](../warehouse-api/README.md))

## Environment Variable

Buat/cek file `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## Menjalankan Aplikasi

```bash
cd warehouse-fe
npm install
npm run dev
```

Buka: `http://localhost:3000`

Default login:

- `admin / admin`

## Build Production

```bash
cd warehouse-fe
npm run build
npm run start
```
