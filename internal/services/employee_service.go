package services

import (
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

// EmployeeService интерфейс для работы с пользователями
type EmployeeService interface {
	Register(user *models.User) error
	GetUserByID(id string) (*models.User, error)
}

type employeeService struct {
	repo repository.EmployeeRepository
}

//NewEmployeeService создает новый экземпляр EmployeeService

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// Register регистрирует нового пользователя
func (s *employeeService) Register(user *models.User) error {
	return s.repo.Create(user)
}

// GetUser ByID получает пользователя по ID
func (s *employeeService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}
