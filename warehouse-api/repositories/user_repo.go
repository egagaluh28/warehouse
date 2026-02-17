package repositories

import (
	"database/sql"
	"warehouse-api/models"
)

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetAll() ([]models.User, error)
	Create(user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, username, email, full_name, role FROM users ORDER BY id DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, password, email, full_name, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.Password, user.Email, user.FullName, user.Role).Scan(&user.ID)
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, email, full_name, role FROM users WHERE username = $1`
	row := r.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, full_name, role FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
