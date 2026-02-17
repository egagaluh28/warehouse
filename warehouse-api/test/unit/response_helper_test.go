package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"warehouse-api/models"
	"warehouse-api/utils"

	"github.com/stretchr/testify/assert"
)

func TestJSONSuccess(t *testing.T) {
	t.Run("Success - Return success response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}
		message := "Operation successful"

		utils.JSONSuccess(w, message, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, message, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("Success - Empty data", func(t *testing.T) {
		w := httptest.NewRecorder()

		utils.JSONSuccess(w, "Success", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
	})
}

func TestJSONCreated(t *testing.T) {
	t.Run("Success - Return 201 created response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"id": "1"}
		message := "Resource created"

		utils.JSONCreated(w, message, data)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response models.APIResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.Equal(t, message, response.Message)
	})
}

func TestJSONError(t *testing.T) {
	t.Run("Error - Return error response with status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		message := "Something went wrong"

		utils.JSONError(w, http.StatusBadRequest, message)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.APIResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.Equal(t, message, response.Message)
	})

	t.Run("Error - Different status codes", func(t *testing.T) {
		testCases := []struct {
			name       string
			statusCode int
			message    string
		}{
			{"Bad Request", http.StatusBadRequest, "Invalid input"},
			{"Unauthorized", http.StatusUnauthorized, "Unauthorized access"},
			{"Forbidden", http.StatusForbidden, "Access denied"},
			{"Not Found", http.StatusNotFound, "Resource not found"},
			{"Internal Server Error", http.StatusInternalServerError, "Server error"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				utils.JSONError(w, tc.statusCode, tc.message)
				assert.Equal(t, tc.statusCode, w.Code)
			})
		}
	})
}

func TestJSONWithMeta(t *testing.T) {
	t.Run("Success - Return response with metadata", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := []string{"item1", "item2"}
		meta := models.Pagination{
			Page:  1,
			Limit: 10,
			Total: 50,
		}

		utils.JSONWithMeta(w, "Data retrieved", data, meta)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Meta)
	})
}
