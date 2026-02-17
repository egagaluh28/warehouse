package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"warehouse-api/models"
	"warehouse-api/repositories"
	"warehouse-api/utils"
)

type BarangHandler struct {
	repo repositories.BarangRepository
}

func NewBarangHandler(repo repositories.BarangRepository) *BarangHandler {
	return &BarangHandler{repo}
}

// GetAll godoc
// @Summary Ambil semua data barang
// @Description Mengambil daftar barang dengan fitur pencarian, pagination, dan sorting.
// @Tags Barang
// @Accept  json
// @Produce  json
// @Param   search query string false "Cari berdasarkan nama/kode"
// @Param   page query int false "Nomor halaman"
// @Param   limit query int false "Jumlah item per halaman"
// @Param   sort_by query string false "Urutkan berdasarkan (harga_beli, harga_jual, nama, id)"
// @Param   order query string false "Urutan (asc, desc)"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /barang [get]
func (h *BarangHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
    sortBy := r.URL.Query().Get("sort_by")
    order := r.URL.Query().Get("order")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	barangs, total, err := h.repo.GetAll(search, limit, offset, sortBy, order)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Server error")
		return
	}

	// Menggunakan helper response standar
	utils.JSONResponse(w, http.StatusOK, true, "Data berhasil diambil", barangs, &models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

// GetAllWithStok godoc
// @Summary Ambil semua data barang dengan stok
// @Description Mengambil daftar barang beserta stok saat ini (join dengan tabel stok) dengan fitur pencarian, pagination, dan sorting.
// @Tags Barang
// @Accept  json
// @Produce  json
// @Param   search query string false "Cari berdasarkan nama/kode"
// @Param   page query int false "Nomor halaman"
// @Param   limit query int false "Jumlah item per halaman"
// @Param   sort_by query string false "Urutkan berdasarkan (harga_beli, harga_jual, nama, id, stok)"
// @Param   order query string false "Urutan (asc, desc)"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /barang/stok [get]
func (h *BarangHandler) GetAllWithStok(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sortBy := r.URL.Query().Get("sort_by")
	order := r.URL.Query().Get("order")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	barangs, total, err := h.repo.GetAllWithStok(search, limit, offset, sortBy, order)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Server error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, "Data berhasil diambil", barangs, &models.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

// GetByID godoc
// @Summary Ambil barang berdasarkan ID
// @Description Mengambil detail barang spesifik
// @Tags Barang
// @Accept  json
// @Produce  json
// @Param   id path int true "ID Barang"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /barang/{id} [get]
func (h *BarangHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	barang, err := h.repo.GetByID(id)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	utils.JSONSuccess(w, "Data retrieved successfully", barang)
}

// Create godoc
// @Summary Tambah barang baru
// @Description Menambahkan data barang baru ke inventaris. kode_barang akan digenerate otomatis oleh sistem.
// @Tags Barang
// @Accept  json
// @Produce  json
// @Param   request body models.CreateBarangRequest true "Data Barang"
// @Security BearerAuth
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /barang [post]
func (h *BarangHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBarangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	if strings.TrimSpace(req.NamaBarang) == "" {
		utils.JSONError(w, http.StatusBadRequest, "Nama barang wajib diisi")
		return
	}
	if req.HargaBeli <= 0 {
		utils.JSONError(w, http.StatusBadRequest, "Harga beli wajib diisi")
		return
	}
	if req.HargaJual <= 0 {
		utils.JSONError(w, http.StatusBadRequest, "Harga jual wajib diisi")
		return
	}
	if strings.TrimSpace(req.Satuan) == "" {
		req.Satuan = "pcs"
	}

	barang := &models.Barang{
		NamaBarang: req.NamaBarang,
		Deskripsi:  req.Deskripsi,
		Satuan:     req.Satuan,
		HargaBeli:  req.HargaBeli,
		HargaJual:  req.HargaJual,
	}

	err := h.repo.Create(barang)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Gagal membuat barang: "+err.Error())
		return
	}

	utils.JSONCreated(w, "Barang berhasil dibuat", barang)
}

// Update godoc
// @Summary Perbarui data barang
// @Description Memperbarui detail barang yang ada
// @Tags Barang
// @Accept  json
// @Produce  json
// @Param   id path int true "ID Barang"
// @Param   request body models.CreateBarangRequest true "Data Barang"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /barang/{id} [put]
func (h *BarangHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	var req models.CreateBarangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	if strings.TrimSpace(req.NamaBarang) == "" {
		utils.JSONError(w, http.StatusBadRequest, "Nama barang wajib diisi")
		return
	}
	if req.HargaBeli <= 0 {
		utils.JSONError(w, http.StatusBadRequest, "Harga beli wajib diisi")
		return
	}
	if req.HargaJual <= 0 {
		utils.JSONError(w, http.StatusBadRequest, "Harga jual wajib diisi")
		return
	}
	if strings.TrimSpace(req.Satuan) == "" {
		req.Satuan = "pcs"
	}

	exists, _ := h.repo.Exists(id)
	if !exists {
		utils.JSONError(w, http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	existing, err := h.repo.GetByID(id)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	barang := &models.Barang{
		ID:         id,
		KodeBarang: existing.KodeBarang,
		NamaBarang: req.NamaBarang,
		Deskripsi:  req.Deskripsi,
		Satuan:     req.Satuan,
		HargaBeli:  req.HargaBeli,
		HargaJual:  req.HargaJual,
	}

	err = h.repo.Update(barang)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Gagal memperbarui barang")
		return
	}

	utils.JSONSuccess(w, "Data barang berhasil diperbarui", barang)
}

// Delete godoc
// @Summary Hapus barang
// @Description Menghapus barang dari inventaris
// @Tags Barang
// @Accept  json
// @Produce  json
// @Param   id path int true "ID Barang"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /barang/{id} [delete]
func (h *BarangHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	exists, _ := h.repo.Exists(id)
	if !exists {
		utils.JSONError(w, http.StatusNotFound, "Barang tidak ditemukan")
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Gagal menghapus barang")
		return
	}

	utils.JSONSuccess(w, "Barang berhasil dihapus", nil)
}
