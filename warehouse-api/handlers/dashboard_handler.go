package handlers

import (
    "net/http"
    "warehouse-api/repositories"
    "warehouse-api/utils"
)

type DashboardHandler struct {
    repo repositories.DashboardRepository
}

func NewDashboardHandler(repo repositories.DashboardRepository) *DashboardHandler {
    return &DashboardHandler{repo}
}

// GetStats godoc
// @Summary Ambil statistik dashboard
// @Description Mengambil ringkasan data gudang (Total Barang, Stok, Nilai Aset, Top Selling)
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /dashboard [get]
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
    stats, err := h.repo.GetStats()
    if err != nil {
        utils.JSONError(w, http.StatusInternalServerError, "Gagal mengambil data dashboard: "+err.Error())
        return
    }

    utils.JSONSuccess(w, "Data dashboard berhasil diambil", stats)
}
