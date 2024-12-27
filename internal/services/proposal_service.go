package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

// ProposalService интерфейс для работы с предложениями
type ProposalService interface {
	Create(proposal *models.Proposal) (int64, error)
	GetByID(id int64) (*models.Proposal, error)
	GetByTenderID(tenderID int64) ([]*models.Proposal, error)
	Update(proposal *models.Proposal) error
	Delete(id int64) error
	List() ([]*models.Proposal, error)
	Publish(id int64) error
	Accept(id int64) error
	Reject(id int64) error
	Cancel(id int64) error
}

type proposalService struct {
	proposalRepo repository.ProposalRepository
	tenderRepo   repository.TenderRepository
}

// NewProposalService создает новый экземпляр ProposalService
func NewProposalService(proposalRepo repository.ProposalRepository, tenderRepo repository.TenderRepository) ProposalService {
	return &proposalService{proposalRepo: proposalRepo, tenderRepo: tenderRepo}
}

// Create создает новое предложение
func (s *proposalService) Create(proposal *models.Proposal) (int64, error) {
	if proposal.Description == "" {
		return 0, errors.New("описание предложения не может быть пустым")
	}

	if proposal.TenderID == 0 {
		return 0, errors.New("необходимо указать тендер")
	}

	if proposal.OrganizationID == 0 {
		return 0, errors.New("необходимо указать организацию")
	}
	proposal.Version = 1
	proposal.CreatedAt = time.Now()
	proposal.PublicationDate = time.Now()

	id, err := s.proposalRepo.Create(proposal)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetByID получает предложение по ID
func (s *proposalService) GetByID(id int64) (*models.Proposal, error) {
	return s.proposalRepo.GetByID(id)
}

// GetByTenderID получает все предложения для определенного тендера
func (s *proposalService) GetByTenderID(tenderID int64) ([]*models.Proposal, error) {
	return s.proposalRepo.GetByTenderID(tenderID)
}

// Update обновляет данные предложения
func (s *proposalService) Update(proposal *models.Proposal) error {
	// Проверка существования предложения
	existingProposal, err := s.proposalRepo.GetByID(proposal.ID)
	if err != nil {
		return err
	}

	// Обновляем версию
	proposal.Version = existingProposal.Version + 1
	proposal.UpdatedAt = time.Now()
	return s.proposalRepo.Update(proposal)
}

// Delete удаляет предложение по ID
func (s *proposalService) Delete(id int64) error {
	return s.proposalRepo.Delete(id)
}

// List возвращает список всех предложений
func (s *proposalService) List() ([]*models.Proposal, error) {
	return s.proposalRepo.List()
}

// Publish устанавливает статус предложения как "published"
func (s *proposalService) Publish(id int64) error {
	proposal, err := s.proposalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("предложение не найдено: %w", err)
	}

	if proposal.Status == models.ProposalStatusPublished {
		return errors.New("предложение уже опубликовано")
	}

	proposal.Status = models.ProposalStatusPublished
	proposal.UpdatedAt = time.Now()
	proposal.PublicationDate = time.Now()
	return s.proposalRepo.Update(proposal)
}

// Accept устанавливает статус предложения как "accepted" и закрывает тендер
func (s *proposalService) Accept(id int64) error {
	proposal, err := s.proposalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("предложение не найдено: %w", err)
	}

	if proposal.Status == models.ProposalStatusAccepted {
		return errors.New("предложение уже принято")
	}

	tender, err := s.tenderRepo.GetByID(proposal.TenderID)
	if err != nil {
		return fmt.Errorf("тендер не найден: %w", err)
	}

	if tender.Status == models.TenderStatusClosed {
		return errors.New("невозможно принять предложение для закрытого тендера")
	}

	if tender.Status == models.TenderStatusCancelled {
		return errors.New("невозможно принять предложение для отмененного тендера")
	}

	proposal.Status = models.ProposalStatusAccepted
	proposal.UpdatedAt = time.Now()
	if err := s.proposalRepo.Update(proposal); err != nil {
		return err
	}

	tender.Status = models.TenderStatusClosed
	tender.UpdatedAt = time.Now()
	return s.tenderRepo.Update(tender)
}

// Reject устанавливает статус предложения как "rejected"
func (s *proposalService) Reject(id int64) error {
	proposal, err := s.proposalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("предложение не найдено: %w", err)
	}

	if proposal.Status == models.ProposalStatusRejected {
		return errors.New("предложение уже отклонено")
	}

	if proposal.Status == models.ProposalStatusAccepted {
		return errors.New("невозможно отклонить принятое предложение")
	}

	proposal.Status = models.ProposalStatusRejected
	proposal.UpdatedAt = time.Now()
	return s.proposalRepo.Update(proposal)
}

// Cancel устанавливает статус предложения как "cancelled"
func (s *proposalService) Cancel(id int64) error {
	proposal, err := s.proposalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("предложение не найдено: %w", err)
	}

	if proposal.Status == models.ProposalStatusCancelled {
		return errors.New("предложение уже отменено")
	}

	if proposal.Status == models.ProposalStatusAccepted {
		return errors.New("невозможно отменить принятое предложение")
	}

	proposal.Status = models.ProposalStatusCancelled
	proposal.UpdatedAt = time.Now()
	return s.proposalRepo.Update(proposal)
}
