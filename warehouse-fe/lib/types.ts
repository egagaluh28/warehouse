// User & Auth
export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  role: "admin" | "staff";
  created_at?: string;
  updated_at?: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  username: string;
  password: string;
  email: string;
  full_name: string;
  role: "admin" | "staff";
}

// Barang (Inventory)
export interface Barang {
  id: number;
  kode_barang: string;
  nama_barang: string;
  deskripsi?: string;
  satuan: string;
  harga_beli: number;
  harga_jual: number;
  stok?: number;
}

export interface CreateBarangRequest {
  nama_barang: string;
  deskripsi: string;
  satuan: string;
  harga_beli: number;
  harga_jual: number;
}

// Stok
export interface Stok {
  id: number;
  barang_id: number;
  stok_akhir: number;
  updated_at: string;
  barang?: Barang;
}

export interface HistoryStok {
  id: number;
  barang_id: number;
  user_id: number;
  jenis_transaksi: string;
  jumlah: number;
  stok_sebelum: number;
  stok_sesudah: number;
  keterangan: string;
  created_at: string;
  barang?: Barang;
  user?: User;
}

// Pembelian (Purchasing)
export interface BeliHeader {
  id: number;
  no_faktur: string;
  supplier: string;
  total: number;
  user_id: number;
  status: string;
  created_at: string;
  user?: User;
  details?: BeliDetail[];
}

export interface BeliDetail {
  id: number;
  beli_header_id: number;
  barang_id: number;
  qty: number;
  harga: number;
  subtotal: number;
  barang?: Barang;
}

export interface CreatePembelianRequest {
  no_faktur?: string;
  supplier: string;
  user_id: number;
  details: {
    barang_id: number;
    qty: number;
    harga: number;
  }[];
}

// Penjualan (Sales)
export interface JualHeader {
  id: number;
  no_faktur: string;
  customer: string;
  total: number;
  user_id: number;
  status: string;
  created_at: string;
  user?: User;
  details?: JualDetail[];
}

export interface JualDetail {
  id: number;
  jual_header_id: number;
  barang_id: number;
  qty: number;
  harga: number;
  subtotal: number;
  barang?: Barang;
}

export interface CreatePenjualanRequest {
  no_faktur?: string;
  customer: string;
  user_id: number;
  details: {
    barang_id: number;
    qty: number;
    harga: number;
  }[];
}

// Dashboard
export interface TopProduct {
  nama_barang: string;
  total_terjual: number;
}

export interface DashboardStats {
  total_user: number;
  total_barang: number;
  total_stok: number;
  total_nilai_aset: number;
  top_selling_products: TopProduct[];
}

// API Utilities
export interface APIResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  message: string;
  data: T[];
  meta: PaginationMeta;
}
