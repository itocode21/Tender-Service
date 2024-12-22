package repository

import (
	"database/sql"
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
)

// EmployeeRepository интерфейс для работы с пользователями
type EmployeeRepository interface {
	Create(user *models.User) (string, error)
	GetByID(id string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
}

// employeeRepository реализация интерфейса EmployeeRepository
type employeeRepository struct {
	db *sql.DB
}

// NewEmployeeRepository создает новый экземпляр employeeRepository
func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

// Create создает нового пользователя
func (r *employeeRepository) Create(user *models.User) (string, error) {
	err := r.db.QueryRow(
		`INSERT INTO employees (username, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`,
		user.Username, user.FirstName, user.LastName).Scan(&user.ID)
	return user.ID, err
}

// GetByID получает пользователя по ID
func (r *employeeRepository) GetByID(id string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT id, username, first_name, last_name, created_at, updated_at FROM employees WHERE id = $1`, id)
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername получает пользователя по имени
func (r *employeeRepository) GetByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT id, username, first_name, last_name, created_at, updated_at FROM employees WHERE username = $1`, username)
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
