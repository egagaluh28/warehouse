import { apiClient } from "./api-client";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  Barang,
  CreateBarangRequest,
  Stok,
  HistoryStok,
  BeliHeader,
  CreatePembelianRequest,
  JualHeader,
  CreatePenjualanRequest,
  DashboardStats,
  APIResponse,
  User,
  PaginatedResponse,
} from "./types";

// Auth API
export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<APIResponse<LoginResponse>>(
      "/login",
      data,
    );
    return response.data.data;
  },

  register: async (data: RegisterRequest): Promise<User> => {
    const response = await apiClient.post<APIResponse<User>>("/register", data);
    return response.data.data;
  },
};

// User API (Admin Only)
export const userApi = {
  getAll: async (): Promise<User[]> => {
    const response = await apiClient.get<APIResponse<User[]>>("/users");
    return response.data.data || [];
  },
};

// Barang API
export const barangApi = {
  getAll: async (params?: {
    search?: string;
    page?: number;
    limit?: number;
    sort_by?: string;
    order?: "asc" | "desc";
  }): Promise<PaginatedResponse<Barang>> => {
    const response = await apiClient.get<PaginatedResponse<Barang>>("/barang", {
      params,
    });
    return response.data;
  },

  getAllWithStok: async (params?: {
    search?: string;
    page?: number;
    limit?: number;
    sort_by?: string;
    order?: "asc" | "desc";
  }): Promise<PaginatedResponse<Barang>> => {
    const response = await apiClient.get<PaginatedResponse<Barang>>(
      "/barang/stok",
      {
        params,
      },
    );
    return response.data;
  },

  getById: async (id: number): Promise<Barang> => {
    const response = await apiClient.get<APIResponse<Barang>>(`/barang/${id}`);
    return response.data.data;
  },

  create: async (data: CreateBarangRequest): Promise<Barang> => {
    const response = await apiClient.post<APIResponse<Barang>>("/barang", data);
    return response.data.data;
  },

  update: async (
    id: number,
    data: Partial<CreateBarangRequest>,
  ): Promise<Barang> => {
    const response = await apiClient.put<APIResponse<Barang>>(
      `/barang/${id}`,
      data,
    );
    return response.data.data;
  },

  delete: async (id: number): Promise<boolean> => {
    const response = await apiClient.delete<APIResponse>(`/barang/${id}`);
    return response.data.success;
  },
};

// Stok API
export const stokApi = {
  getAll: async (): Promise<Stok[]> => {
    const response = await apiClient.get<APIResponse<Stok[]>>("/stok");
    return response.data.data || [];
  },

  // Note: Backend might need an endpoint for single stock if not available
  // Assuming /stok/{id} or list filter
  getHistory: async (barangId?: number): Promise<HistoryStok[]> => {
    const url = barangId ? `/history-stok/${barangId}` : "/history-stok";
    const response = await apiClient.get<APIResponse<HistoryStok[]>>(url);
    return response.data.data || [];
  },

  getById: async (id: number): Promise<Stok> => {
    const response = await apiClient.get<APIResponse<Stok>>(`/stok/${id}`);
    return response.data.data;
  },
};

// Pembelian API
export const pembelianApi = {
  getAll: async (params?: {
    start_date?: string;
    end_date?: string;
  }): Promise<BeliHeader[]> => {
    const response = await apiClient.get<APIResponse<BeliHeader[]>>(
      "/pembelian",
      { params },
    );
    return response.data.data || [];
  },

  getById: async (id: number): Promise<BeliHeader> => {
    const response = await apiClient.get<APIResponse<BeliHeader>>(
      `/pembelian/${id}`,
    );
    return response.data.data;
  },

  create: async (data: CreatePembelianRequest): Promise<BeliHeader> => {
    const response = await apiClient.post<APIResponse<BeliHeader>>(
      "/pembelian",
      data,
    );
    return response.data.data;
  },
};

// Penjualan API
export const penjualanApi = {
  getAll: async (params?: {
    start_date?: string;
    end_date?: string;
  }): Promise<JualHeader[]> => {
    const response = await apiClient.get<APIResponse<JualHeader[]>>(
      "/penjualan",
      { params },
    );
    return response.data.data || [];
  },

  getById: async (id: number): Promise<JualHeader> => {
    const response = await apiClient.get<APIResponse<JualHeader>>(
      `/penjualan/${id}`,
    );
    return response.data.data;
  },

  create: async (data: CreatePenjualanRequest): Promise<JualHeader> => {
    const response = await apiClient.post<APIResponse<JualHeader>>(
      "/penjualan",
      data,
    );
    return response.data.data;
  },
};

// Dashboard API
export const dashboardApi = {
  getStats: async (): Promise<DashboardStats> => {
    const response =
      await apiClient.get<APIResponse<DashboardStats>>("/dashboard");
    return response.data.data;
  },
};
