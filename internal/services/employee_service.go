package services

import (
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

// EmployeeService интерфейс для работы с пользователями
type EmployeeService interface {
	Register(user *models.User) (int, error)
	GetUserByID(id int) (*models.User, error)
}

type employeeService struct {
	repo repository.EmployeeRepository
}

// NewEmployeeService создает новый экземпляр EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// Register регистрирует нового пользователя
func (s *employeeService) Register(user *models.User) (int, error) {
	if user.Username == "" {
		return 0, errors.New("Имя пользователя не может быть пустым") // Возвращаем 0 вместо "" для int
	}
	id, err := s.repo.Create(user) // Получаем ID и ошибку
	return id, err                 // Возвращаем оба значения
}

// GetUser ByID получает пользователя по ID
func (s *employeeService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}
