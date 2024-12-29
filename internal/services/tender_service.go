package services

import (
	"errors"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

type TenderService interface {
	Create(tender *models.Tender) (int64, error)
	GetByID(id int64) (*models.Tender, error)
	Update(tender *models.Tender) error
	Delete(id int64) error
	List() ([]*models.Tender, error)
	ListByOrganizationID(organizationId int64) ([]*models.Tender, error)
	Publish(id int64) error
	Close(id int64) error
}

type tenderService struct {
	tenderRepo repository.TenderRepository
}

func NewTenderService(tenderRepo repository.TenderRepository) TenderService {
	return &tenderService{tenderRepo: tenderRepo}
}

func (s *tenderService) Create(tender *models.Tender) (int64, error) {
	if tender.Name == "" {
		return 0, errors.New("имя тендера не может быть пустым")
	}
	if tender.OrganizationID == 0 {
		return 0, errors.New("идентификатор организации не может быть пустым")
	}
	tender.Version = 1
	tender.CreatedAt = time.Now()

	id, err := s.tenderRepo.Create(tender)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *tenderService) GetByID(id int64) (*models.Tender, error) {
	return s.tenderRepo.GetByID(id)
}

func (s *tenderService) Update(tender *models.Tender) error {
	if tender.Name == "" {
		return errors.New("имя тендера не может быть пустым")
	}

	existingTender, err := s.tenderRepo.GetByID(tender.ID)
	if err != nil {
		return err
	}

	tender.Version = existingTender.Version + 1
	tender.UpdatedAt = time.Now()

	return s.tenderRepo.Update(tender)
}

func (s *tenderService) Delete(id int64) error {
	return s.tenderRepo.Delete(id)
}

func (s *tenderService) List() ([]*models.Tender, error) {
	return s.tenderRepo.List()
}

func (s *tenderService) ListByOrganizationID(organizationId int64) ([]*models.Tender, error) {
	return s.tenderRepo.ListByOrganizationID(organizationId)
}

func (s *tenderService) Publish(id int64) error {
	tender, err := s.tenderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if tender.Status == models.TenderStatusPublished {
		return errors.New("тендер уже опубликован")
	}

	tender.Status = models.TenderStatusPublished
	tender.UpdatedAt = time.Now()

	return s.tenderRepo.Update(tender)
}

func (s *tenderService) Close(id int64) error {
	tender, err := s.tenderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if tender.Status == models.TenderStatusClosed {
		return errors.New("тендер уже закрыт")
	}

	if tender.Status == models.TenderStatusCancelled {
		return errors.New("невозможно закрыть отмененный тендер")
	}

	tender.Status = models.TenderStatusClosed
	tender.UpdatedAt = time.Now()
	return s.tenderRepo.Update(tender)
}
