package handlers

import (
	"net/http"
	"strconv"
	
	"warehouse-api/repositories"
    "warehouse-api/utils"
)

type StokHandler struct {
	repo repositories.StokRepository
}

func NewStokHandler(repo repositories.StokRepository) *StokHandler {
	return &StokHandler{repo}
}


// GetAll godoc
// @Summary Ambil semua stok
// @Description Mengambil data stok terkini untuk semua barang
// @Tags Stok
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /stok [get]
func (h *StokHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	stoks, err := h.repo.GetAll()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Server error")
		return
	}

	utils.JSONSuccess(w, "Data stok berhasil diambil", stoks)
}

// GetByBarangID godoc
// @Summary Ambil stok berdasarkan ID barang
// @Description Mengambil informasi stok untuk barang tertentu
// @Tags Stok
// @Accept  json
// @Produce  json
// @Param   id path int true "ID Barang"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /stok/{id} [get]
func (h *StokHandler) GetByBarangID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	stok, err := h.repo.GetByBarangID(id)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "Stok tidak ditemukan")
		return
	}

	utils.JSONSuccess(w, "Data stok berhasil diambil", stok)
}

// GetHistory godoc
// @Summary Ambil riwayat stok
// @Description Mengambil riwayat pergerakan stok. Opsional: filter berdasarkan ID barang.
// @Tags Stok
// @Accept  json
// @Produce  json
// @Param   id path int false "ID Barang (opsional)"
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /history-stok [get]
// @Router /history-stok/{id} [get]
func (h *StokHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id") // Opsional, untuk riwayat barang tertentu
    var barangID int
    if idStr != "" {
        var err error
        barangID, err = strconv.Atoi(idStr)
        if err != nil {
            utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
            return
        }
    }

	history, err := h.repo.GetHistory(barangID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Server error")
		return
	}

	utils.JSONSuccess(w, "Riwayat stok berhasil diambil", history)
}
