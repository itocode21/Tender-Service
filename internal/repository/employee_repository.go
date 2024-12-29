package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/itocode21/Tender-Service/internal/models"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type EmployeeRepository interface {
	Create(employee *models.User) (int64, error)
	GetByID(id int64) (*models.User, error)
	Update(employee *models.User) error
	Delete(id int64) error
	List() ([]*models.User, error)
}

type EmployeeRepositoryImpl struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{db: db}
}

func (r *EmployeeRepositoryImpl) Create(employee *models.User) (int64, error) {
	query := `INSERT INTO employees (username, first_name, last_name, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64

	err := r.db.QueryRow(query, employee.Username, employee.FirstName, employee.LastName, employee.CreatedAt, employee.UpdatedAt).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка при создании сотрудника: %w", err)
	}

	return id, nil
}

func (r *EmployeeRepositoryImpl) GetByID(id int64) (*models.User, error) {
	query := `SELECT id, username, first_name, last_name, created_at, updated_at FROM employees WHERE id = $1`
	var employee models.User

	err := r.db.QueryRow(query, id).Scan(&employee.Id, &employee.Username, &employee.FirstName, &employee.LastName, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("employee not found")
		}
		return nil, fmt.Errorf("ошибка при получении сотрудника по ID: %w", err)
	}

	return &employee, nil
}

func (r *EmployeeRepositoryImpl) Update(employee *models.User) error {
	query := `UPDATE employees SET username = $1, first_name = $2, last_name = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.Exec(query, employee.Username, employee.FirstName, employee.LastName, employee.UpdatedAt, employee.Id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении сотрудника: %w", err)
	}
	return nil
}

func (r *EmployeeRepositoryImpl) Delete(id int64) error {
	query := `DELETE FROM employees WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении сотрудника: %w", err)
	}
	return nil
}

func (r *EmployeeRepositoryImpl) List() ([]*models.User, error) {
	query := `SELECT id, username, first_name, last_name, created_at, updated_at FROM employees`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка сотрудников: %w", err)
	}
	defer rows.Close()

	var employees []*models.User
	for rows.Next() {
		var employee models.User
		if err := rows.Scan(&employee.Id, &employee.Username, &employee.FirstName, &employee.LastName, &employee.CreatedAt, &employee.UpdatedAt); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки сотрудника: %w", err)
		}
		employees = append(employees, &employee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после итерации по строкам: %w", err)
	}

	return employees, nil
}
