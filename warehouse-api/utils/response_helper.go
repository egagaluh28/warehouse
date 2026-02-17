package utils

import (
	"encoding/json"
	"net/http"
	"warehouse-api/models"
)

// JSONResponse sends a standardized JSON response
func JSONResponse(w http.ResponseWriter, code int, success bool, message string, data interface{}, meta *models.Pagination) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := models.APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// JSONError sends a standardized error response
func JSONError(w http.ResponseWriter, code int, message string) {
	JSONResponse(w, code, false, message, nil, nil)
}

// JSONSuccess sends a standardized success response
func JSONSuccess(w http.ResponseWriter, message string, data interface{}) {
	JSONResponse(w, http.StatusOK, true, message, data, nil)
}

// JSONCreated sends a specialized success response for resource creation
func JSONCreated(w http.ResponseWriter, message string, data interface{}) {
	JSONResponse(w, http.StatusCreated, true, message, data, nil)
}

// JSONWithMeta sends a success response with pagination metadata
func JSONWithMeta(w http.ResponseWriter, message string, data interface{}, meta models.Pagination) {
	JSONResponse(w, http.StatusOK, true, message, data, &meta)
}
