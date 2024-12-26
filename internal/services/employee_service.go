package services

import (
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

// EmployeeService интерфейс для работы с пользователями
type EmployeeService interface {
	Create(user *models.User) (int, error)
	GetByID(id int) (*models.User, error)
	Update(employee *models.User) error
	Delete(id int) error
	List() ([]*models.User, error)
}

type employeeService struct {
	repo repository.EmployeeRepository
}

// NewEmployeeService создает новый экземпляр EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// Create создает нового пользователя
func (s *employeeService) Create(user *models.User) (int, error) {
	if user.Username == "" {
		return 0, errors.New("имя пользователя не может быть пустым")
	}
	id, err := s.repo.Create(user)
	return int(id), err
}

// GetByID получает пользователя по ID
func (s *employeeService) GetByID(id int) (*models.User, error) {
	emp, err := s.repo.GetByID(int64(id))
	return emp, err
}

// Update обновляет данные сотрудника
func (s *employeeService) Update(employee *models.User) error {
	return s.repo.Update(employee)
}

// Delete удаляет сотрудника по ID
func (s *employeeService) Delete(id int) error {
	return s.repo.Delete(int64(id))
}

// List возвращает список всех сотрудников
func (s *employeeService) List() ([]*models.User, error) {
	return s.repo.List()
}
