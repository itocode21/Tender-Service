package services

import (
	"errors"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

// TenderService интерфейс для работы с тендерами
type TenderService interface {
	Create(tender *models.Tender) (int64, error)
	GetByID(id int64) (*models.Tender, error)
	GetByOrganizationID(organizationID int64) ([]*models.Tender, error)
	Update(tender *models.Tender) error
	Delete(id int64) error
	List() ([]*models.Tender, error)
	Publish(id int64) error
	Close(id int64) error
	Cancel(id int64) error
}

type tenderService struct {
	tenderRepo repository.TenderRepository
}

// NewTenderService создает новый экземпляр TenderService
func NewTenderService(tenderRepo repository.TenderRepository) TenderService {
	return &tenderService{tenderRepo: tenderRepo}
}

// Create создает новый тендер
func (s *tenderService) Create(tender *models.Tender) (int64, error) {
	if tender.Name == "" {
		return 0, errors.New("название тендера не может быть пустым")
	}

	if tender.OrganizationID == 0 {
		return 0, errors.New("необходимо указать организацию")
	}

	tender.Status = models.TenderStatusDraft
	tender.PublicationDate = time.Now() // Устанавливаем время публикации
	id, err := s.tenderRepo.Create(tender)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetByID получает тендер по ID
func (s *tenderService) GetByID(id int64) (*models.Tender, error) {
	return s.tenderRepo.GetByID(id)
}

// GetByOrganizationID получает все тендеры по ID организации
func (s *tenderService) GetByOrganizationID(organizationID int64) ([]*models.Tender, error) {
	return s.tenderRepo.GetByOrganizationID(organizationID)
}

// Update обновляет данные тендера
func (s *tenderService) Update(tender *models.Tender) error {
	return s.tenderRepo.Update(tender)
}

// Delete удаляет тендер по ID
func (s *tenderService) Delete(id int64) error {
	return s.tenderRepo.Delete(id)
}

// List возвращает список всех тендеров
func (s *tenderService) List() ([]*models.Tender, error) {
	return s.tenderRepo.List()
}

// Publish устанавливает статус тендера как "published"
func (s *tenderService) Publish(id int64) error {
	tender, err := s.tenderRepo.GetByID(id)
	if err != nil {
		return err
	}
	if tender.Status == models.TenderStatusPublished {
		return errors.New("тендер уже опубликован")
	}

	tender.Status = models.TenderStatusPublished
	tender.PublicationDate = time.Now()
	return s.tenderRepo.Update(tender)
}

// Close устанавливает статус тендера как "closed"
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
	return s.tenderRepo.Update(tender)
}

// Cancel устанавливает статус тендера как "cancelled"
func (s *tenderService) Cancel(id int64) error {
	tender, err := s.tenderRepo.GetByID(id)
	if err != nil {
		return err
	}
	if tender.Status == models.TenderStatusCancelled {
		return errors.New("тендер уже отменен")
	}
	if tender.Status == models.TenderStatusClosed {
		return errors.New("невозможно отменить закрытый тендер")
	}
	tender.Status = models.TenderStatusCancelled
	return s.tenderRepo.Update(tender)
}
