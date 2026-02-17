package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"warehouse-api/handlers"
	"warehouse-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Barang Repository for Handler Tests
type MockBarangRepositoryHandler struct {
	mock.Mock
}

func (m *MockBarangRepositoryHandler) GetAll(search string, limit, offset int, sortBy, order string) ([]models.Barang, int, error) {
	args := m.Called(search, limit, offset, sortBy, order)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Barang), args.Int(1), args.Error(2)
}

func (m *MockBarangRepositoryHandler) GetAllWithStok(search string, limit, offset int, sortBy, order string) ([]models.BarangWithStok, int, error) {
	args := m.Called(search, limit, offset, sortBy, order)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.BarangWithStok), args.Int(1), args.Error(2)
}

func (m *MockBarangRepositoryHandler) GetByID(id int) (*models.BarangWithStok, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BarangWithStok), args.Error(1)
}

func (m *MockBarangRepositoryHandler) GetByKode(kode string) (*models.Barang, error) {
	args := m.Called(kode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Barang), args.Error(1)
}

func (m *MockBarangRepositoryHandler) Create(barang *models.Barang) error {
	args := m.Called(barang)
	return args.Error(0)
}

func (m *MockBarangRepositoryHandler) Update(barang *models.Barang) error {
	args := m.Called(barang)
	return args.Error(0)
}

func (m *MockBarangRepositoryHandler) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBarangRepositoryHandler) Exists(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestBarangHandlerGetAll(t *testing.T) {
	t.Run("Success - Get all barang with pagination", func(t *testing.T) {
		mockRepo := new(MockBarangRepositoryHandler)
		handler := handlers.NewBarangHandler(mockRepo)

		expectedBarang := []models.Barang{
			{ID: 1, KodeBarang: "BRG-001", NamaBarang: "Item 1", Satuan: "pcs", HargaBeli: 1000, HargaJual: 1500},
			{ID: 2, KodeBarang: "BRG-002", NamaBarang: "Item 2", Satuan: "pcs", HargaBeli: 2000, HargaJual: 2500},
		}

		mockRepo.On("GetAll", "", 10, 0, "", "").Return(expectedBarang, 2, nil)

		req := httptest.NewRequest("GET", "/api/barang?page=1&limit=10", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Search barang", func(t *testing.T) {
		mockRepo := new(MockBarangRepositoryHandler)
		handler := handlers.NewBarangHandler(mockRepo)

		expectedBarang := []models.Barang{
			{ID: 1, KodeBarang: "BRG-001", NamaBarang: "Laptop", Satuan: "unit", HargaBeli: 5000000, HargaJual: 6000000},
		}

		mockRepo.On("GetAll", "Laptop", 10, 0, "", "").Return(expectedBarang, 1, nil)

		req := httptest.NewRequest("GET", "/api/barang?search=Laptop&page=1&limit=10", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBarangHandlerGetByID(t *testing.T) {
	t.Run("Success - Get barang by ID", func(t *testing.T) {
		mockRepo := new(MockBarangRepositoryHandler)
		handler := handlers.NewBarangHandler(mockRepo)

		expectedBarang := &models.BarangWithStok{
			Barang: models.Barang{
				ID:         1,
				KodeBarang: "BRG-001",
				NamaBarang: "Test Item",
				Satuan:     "pcs",
				HargaBeli:  1000,
				HargaJual:  1500,
			},
			Stok: 100,
		}

		mockRepo.On("GetByID", 1).Return(expectedBarang, nil)

		req := httptest.NewRequest("GET", "/api/barang/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Invalid ID format", func(t *testing.T) {
		mockRepo := new(MockBarangRepositoryHandler)
		handler := handlers.NewBarangHandler(mockRepo)

		req := httptest.NewRequest("GET", "/api/barang/invalid", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
