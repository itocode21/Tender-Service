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

// NewTenderService создает новый экземпляр TenderService
func NewTenderService(tenderRepo repository.TenderRepository) TenderService {
	return &tenderService{tenderRepo: tenderRepo}
}

// Create создает новый тендер
func (s *tenderService) Create(tender *models.Tender) (int64, error) {
	// Проверка имени
	if tender.Name == "" {
		return 0, errors.New("имя тендера не может быть пустым")
	}
	// Проверка организации
	if tender.OrganizationID == 0 {
		return 0, errors.New("идентификатор организации не может быть пустым")
	}
	// Устанавливаем начальную версию
	tender.Version = 1
	// Устанавливаем время создания
	tender.CreatedAt = time.Now()

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

// Update обновляет данные тендера
func (s *tenderService) Update(tender *models.Tender) error {
	if tender.Name == "" {
		return errors.New("имя тендера не может быть пустым")
	}

	// Проверка существования тендера
	existingTender, err := s.tenderRepo.GetByID(tender.ID)
	if err != nil {
		return err
	}

	// Обновляем версию
	tender.Version = existingTender.Version + 1
	tender.UpdatedAt = time.Now()

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

// ListByOrganizationID возвращает список тендеров по ID организации
func (s *tenderService) ListByOrganizationID(organizationId int64) ([]*models.Tender, error) {
	return s.tenderRepo.ListByOrganizationID(organizationId)
}

// Publish публикует тендер по ID
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

// Close закрывает тендер по ID
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
