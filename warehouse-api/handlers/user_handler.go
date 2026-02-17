package handlers

import (
    "encoding/json"
    "net/http"
    "os"
    "time"
    "warehouse-api/models"
    "warehouse-api/services"
    
    "warehouse-api/middleware"
    "warehouse-api/utils"
    "github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
    service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
    return &UserHandler{service}
}

// Register godoc
// @Summary Mendaftarkan pengguna baru
// @Description Mendaftarkan pengguna baru (hanya admin)
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body models.RegisterRequest true "Informasi Pengguna"
// @Security BearerAuth
// @Success 201 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Cek Peran Pengguna
    role := r.Context().Value(middleware.RoleKey).(string)
    if role != "admin" {
         utils.JSONError(w, http.StatusForbidden, "Akses ditolak: Hanya admin yang dapat mendaftarkan pengguna baru")
         return
    }

    // 2. Parsing Data Masukan
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
         utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
         return
    }

    user, err := h.service.Register(&req)
    if err != nil {
         utils.JSONError(w, http.StatusBadRequest, err.Error())
         return
    }

    utils.JSONCreated(w, "Pengguna berhasil didaftarkan", user)
}

// GetAll Users
// @Summary Mendapatkan semua pengguna
// @Description Mendapatkan daftar semua pengguna (hanya admin)
// @Tags Auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    // 1. Cek Peran Pengguna
    role := r.Context().Value(middleware.RoleKey).(string)
    if role != "admin" {
         utils.JSONError(w, http.StatusForbidden, "Akses ditolak: Hanya admin yang dapat melihat daftar pengguna")
         return
    }

    // 2. Ambil Data
    users, err := h.service.GetAll()
    if err != nil {
         utils.JSONError(w, http.StatusInternalServerError, "Gagal mengambil data pengguna")
         return
    }

    utils.JSONSuccess(w, "Berhasil mengambil data pengguna", users)
}

// Login godoc
// @Summary Masuk sistem
// @Description Otentikasi pengguna dan mendapatkan token JWT
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body models.LoginRequest true "Kredensial Login"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.JSONError(w, http.StatusBadRequest, "Data input tidak valid")
        return
    }

    // Validate credentials using service
    user, err := h.service.ValidateCredentials(req.Username, req.Password)
    if err != nil {
        utils.JSONError(w, http.StatusUnauthorized, err.Error())
        return
    }

    // Generate Token JWT
    secret := []byte(os.Getenv("JWT_SECRET"))
    if len(secret) == 0 {
        secret = []byte("supersecretkey")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "username": user.Username,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(secret)
    if err != nil {
        utils.JSONError(w, http.StatusInternalServerError, "Gagal membuat token")
        return
    }

    resp := models.APIResponse{
        Success: true,
        Message: "Login berhasil",
        Data: models.LoginResponse{
            Token: tokenString,
            User:  *user,
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
