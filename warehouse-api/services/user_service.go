package services

import (
	"errors"
	"warehouse-api/models"
	"warehouse-api/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	ValidateCredentials(username, password string) (*models.User, error)
	GetAll() ([]models.User, error)
	HashPassword(password string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

// Register handles user registration with password hashing
func (s *userService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username dan password harus diisi")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password minimal 6 karakter")
	}

	if req.Role != "admin" && req.Role != "staff" {
		return nil, errors.New("role harus 'admin' atau 'staff'")
	}

	// Hash password
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user object
	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     req.Role,
	}

	// Save to database
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateCredentials checks if username and password are correct
func (s *userService) ValidateCredentials(username, password string) (*models.User, error) {
	// Get user from database
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	return user, nil
}

// HashPassword generates bcrypt hash from plain password
func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("gagal mengenkripsi password")
	}
	return string(bytes), nil
}
