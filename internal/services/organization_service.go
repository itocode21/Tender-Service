package services

import (
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

// OrganizationService интерфейс для работы с организациями
type OrganizationService interface {
	Create(org *models.Organization) (int64, error)
	GetByID(id int64) (*models.Organization, error)
	Update(org *models.Organization) error
	Delete(id int64) error
	List() ([]*models.Organization, error)
	AddResponsible(orgResp *models.OrganizationResponsible) (int64, error)
	RemoveResponsible(orgResp *models.OrganizationResponsible) error
}

type organizationService struct {
	orgRepo repository.OrganizationRepository
}

// NewOrganizationService создает новый экземпляр OrganizationService
func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{orgRepo: orgRepo}
}

// Create создает новую организацию
func (s *organizationService) Create(org *models.Organization) (int64, error) {
	if org.Name == "" {
		return 0, errors.New("имя организации не может быть пустым")
	}

	id, err := s.orgRepo.Create(org)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetByID получает организацию по ID
func (s *organizationService) GetByID(id int64) (*models.Organization, error) {
	return s.orgRepo.GetByID(id)
}

// Update обновляет данные организации
func (s *organizationService) Update(org *models.Organization) error {
	if org.Name == "" {
		return errors.New("имя организации не может быть пустым")
	}
	return s.orgRepo.Update(org)
}

// Delete удаляет организацию по ID
func (s *organizationService) Delete(id int64) error {
	return s.orgRepo.Delete(id)
}

// List возвращает список всех организаций
func (s *organizationService) List() ([]*models.Organization, error) {
	return s.orgRepo.List()
}

// AddResponsible добавляет ответственного за организацию
func (s *organizationService) AddResponsible(orgResp *models.OrganizationResponsible) (int64, error) {
	if orgResp.OrganizationID == nil || orgResp.UserID == nil {
		return 0, errors.New("организация и пользователь не могут быть nil")
	}
	return s.orgRepo.AddResponsible(orgResp)
}

// RemoveResponsible удаляет ответственного за организацию
func (s *organizationService) RemoveResponsible(orgResp *models.OrganizationResponsible) error {
	return s.orgRepo.RemoveResponsible(orgResp)
}
