package unit

import (
	"errors"
	"testing"
	"warehouse-api/models"
	"warehouse-api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Test HashPassword
func TestHashPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	t.Run("Success - Hash valid password", func(t *testing.T) {
		password := "mySecurePass123"
		hash, err := service.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
		assert.Greater(t, len(hash), 50) // bcrypt hash is typically 60 chars
	})

	t.Run("Success - Different passwords produce different hashes", func(t *testing.T) {
		hash1, _ := service.HashPassword("password1")
		hash2, _ := service.HashPassword("password2")

		assert.NotEqual(t, hash1, hash2)
	})
}

// Test Register
func TestRegister(t *testing.T) {
	t.Run("Success - Register valid user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		req := &models.RegisterRequest{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     "staff",
		}

		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, err := service.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.NotEqual(t, req.Password, user.Password) // Password should be hashed
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Empty username", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		req := &models.RegisterRequest{
			Username: "",
			Password: "password123",
			Email:    "test@example.com",
			Role:     "staff",
		}

		user, err := service.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "username dan password")
	})

	t.Run("Fail - Empty password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		req := &models.RegisterRequest{
			Username: "testuser",
			Password: "",
			Email:    "test@example.com",
			Role:     "staff",
		}

		user, err := service.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("Fail - Password too short", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		req := &models.RegisterRequest{
			Username: "testuser",
			Password: "12345", // Less than 6 characters
			Email:    "test@example.com",
			Role:     "staff",
		}

		user, err := service.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "password minimal 6")
	})

	t.Run("Fail - Invalid role", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		req := &models.RegisterRequest{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
			Role:     "invalidrole",
		}

		user, err := service.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "role harus")
	})

	t.Run("Fail - Database error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		req := &models.RegisterRequest{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
			Role:     "staff",
		}

		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("database error"))

		user, err := service.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

// Test ValidateCredentials
func TestValidateCredentials(t *testing.T) {
	t.Run("Success - Valid credentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		password := "password123"
		hashedPassword, _ := service.HashPassword(password)

		expectedUser := &models.User{
			ID:       1,
			Username: "testuser",
			Password: hashedPassword,
			Role:     "staff",
		}

		mockRepo.On("GetByUsername", "testuser").Return(expectedUser, nil)

		user, err := service.ValidateCredentials("testuser", password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Username, user.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - User not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("user not found"))

		user, err := service.ValidateCredentials("nonexistent", "password123")

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)

		correctPassword := "password123"
		hashedPassword, _ := service.HashPassword(correctPassword)

		existingUser := &models.User{
			ID:       1,
			Username: "testuser",
			Password: hashedPassword,
		}

		mockRepo.On("GetByUsername", "testuser").Return(existingUser, nil)

		user, err := service.ValidateCredentials("testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "username atau password")
		mockRepo.AssertExpectations(t)
	})
}
