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

type PenjualanHandler struct {
    service services.PenjualanService
	repo    repositories.PenjualanRepository
}

func NewPenjualanHandler(service services.PenjualanService, repo repositories.PenjualanRepository) *PenjualanHandler {
	return &PenjualanHandler{service, repo}
}

// Create godoc
// @Summary Buat transaksi penjualan
// @Description Mencatat transaksi penjualan baru dan mengurangi stok
// @Tags Penjualan
// @Accept  json
// @Produce  json
// @Param   request body models.CreatePenjualanRequest true "Data Penjualan"
// @Security BearerAuth
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /penjualan [post]
func (h *PenjualanHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req models.CreatePenjualanRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
        return
    }

    userID := r.Context().Value(middleware.UserIDKey).(int)
    req.UserID = userID

    header, err := h.service.Create(req)
    if err != nil {
        // Cek error string untuk menentukan status code (sederhana)
        // Idealnya menggunakan error type khusus
        if err.Error()[:4] == "stok" || err.Error()[:6] == "barang" {
             utils.JSONError(w, http.StatusBadRequest, err.Error())
        } else {
             utils.JSONError(w, http.StatusInternalServerError, "Gagal memproses transaksi: "+err.Error())
        }
        return
    }

    utils.JSONCreated(w, "Penjualan berhasil dibuat", header)
}

// GetAll godoc
// @Summary Ambil semua transaksi penjualan
// @Description Mengambil daftar transaksi penjualan. Mendukung filter rentang tanggal.
// @Tags Penjualan
// @Accept  json
// @Produce  json
// @Param   start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param   end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /penjualan [get]
func (h *PenjualanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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
// @Summary Ambil detail penjualan
// @Description Mengambil detail transaksi penjualan spesifik
// @Tags Penjualan
// @Accept  json
// @Produce  json
// @Param   id path int true "ID Transaksi"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /penjualan/{id} [get]
func (h *PenjualanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
