package services

import (
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

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

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{orgRepo: orgRepo}
}

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

func (s *organizationService) GetByID(id int64) (*models.Organization, error) {
	return s.orgRepo.GetByID(id)
}

func (s *organizationService) Update(org *models.Organization) error {
	if org.Name == "" {
		return errors.New("имя организации не может быть пустым")
	}
	return s.orgRepo.Update(org)
}

func (s *organizationService) Delete(id int64) error {
	return s.orgRepo.Delete(id)
}

func (s *organizationService) List() ([]*models.Organization, error) {
	return s.orgRepo.List()
}

func (s *organizationService) AddResponsible(orgResp *models.OrganizationResponsible) (int64, error) {
	if orgResp.OrganizationID == nil || orgResp.UserID == nil {
		return 0, errors.New("организация и пользователь не могут быть nil")
	}
	return s.orgRepo.AddResponsible(orgResp)
}

func (s *organizationService) RemoveResponsible(orgResp *models.OrganizationResponsible) error {
	return s.orgRepo.RemoveResponsible(orgResp)
}
