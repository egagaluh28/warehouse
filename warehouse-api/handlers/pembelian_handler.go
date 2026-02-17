package handlers

import (
	"encoding/json"
	"net/http"
    "strconv"
    // "fmt" // Removed unused import
	"warehouse-api/models"
	"warehouse-api/repositories"
    "warehouse-api/services"
    "warehouse-api/middleware"
    "warehouse-api/utils"
)

type PembelianHandler struct {
    service services.PembelianService
	repo    repositories.PembelianRepository
}

func NewPembelianHandler(service services.PembelianService, repo repositories.PembelianRepository) *PembelianHandler {
	return &PembelianHandler{service, repo}
}

// Create godoc
// @Summary Buat transaksi pembelian
// @Description Mencatat transaksi pembelian stok masuk baru
// @Tags Pembelian
// @Accept  json
// @Produce  json
// @Param   request body models.CreatePembelianRequest true "Data Pembelian"
// @Security BearerAuth
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /pembelian [post]
func (h *PembelianHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req models.CreatePembelianRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
        return
    }

    userID := r.Context().Value(middleware.UserIDKey).(int)
    req.UserID = userID

    header, err := h.service.Create(req)
    if err != nil {
        utils.JSONError(w, http.StatusInternalServerError, "Gagal membuat transaksi: "+err.Error())
        return
    }

	utils.JSONCreated(w, "Pembelian berhasil dibuat", header)
}

// GetAll godoc
// @Summary Ambil semua transaksi pembelian
// @Description Mengambil daftar transaksi pembelian. Mendukung filter rentang tanggal.
// @Tags Pembelian
// @Accept  json
// @Produce  json
// @Param   start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param   end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /pembelian [get]
func (h *PembelianHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    startDate := r.URL.Query().Get("start_date")
    endDate := r.URL.Query().Get("end_date")

    transaksi, err := h.repo.GetAll(startDate, endDate)
    if err != nil {
        utils.JSONError(w, http.StatusInternalServerError, "Server error")
        return
    }

    utils.JSONSuccess(w, "Data berhasil diambil", transaksi)
}

// GetByID godoc
// @Summary Ambil detail pembelian
// @Description Mengambil detail transaksi pembelian spesifik
// @Tags Pembelian
// @Accept  json
// @Produce  json
// @Param   id path int true "ID Transaksi"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /pembelian/{id} [get]
func (h *PembelianHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
        return
    }

    transaksi, err := h.repo.GetByID(id)
    if err != nil {
        utils.JSONError(w, http.StatusNotFound, "Transaksi tidak ditemukan")
        return
    }

    utils.JSONSuccess(w, "Data berhasil diambil", transaksi)
}
