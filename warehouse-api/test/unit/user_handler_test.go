package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"warehouse-api/handlers"
	"warehouse-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock User Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(req *models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ValidateCredentials(username, password string) (*models.User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func TestUserHandlerLogin(t *testing.T) {
	t.Run("Success - Valid login", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.NewUserHandler(mockService)

		loginReq := models.LoginRequest{
			Username: "admin",
			Password: "admin123",
		}

		user := &models.User{
			ID:       1,
			Username: "admin",
			Role:     "admin",
		}

		mockService.On("ValidateCredentials", loginReq.Username, loginReq.Password).Return(user, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		mockService.AssertExpectations(t)
	})

	t.Run("Fail - Invalid JSON", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.NewUserHandler(mockService)

		req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandlerRegister(t *testing.T) {
	t.Run("Requires auth middleware context", func(t *testing.T) {
		// Register endpoint requires auth middleware to check user role
		// This would need complex mocking of context values
		// Integration tests cover this scenario better
		t.Skip("Requires auth middleware context - tested in integration tests")
	})
}
