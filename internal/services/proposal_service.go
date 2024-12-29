package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
)

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

func NewProposalService(proposalRepo repository.ProposalRepository, tenderRepo repository.TenderRepository) ProposalService {
	return &proposalService{proposalRepo: proposalRepo, tenderRepo: tenderRepo}
}

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

func (s *proposalService) GetByID(id int64) (*models.Proposal, error) {
	return s.proposalRepo.GetByID(id)
}

func (s *proposalService) GetByTenderID(tenderID int64) ([]*models.Proposal, error) {
	_, err := s.tenderRepo.GetByID(tenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tender not found")
		}
		return nil, fmt.Errorf("failed to get tender: %w", err)
	}
	return s.proposalRepo.GetByTenderID(tenderID)
}

func (s *proposalService) Update(proposal *models.Proposal) error {
	existingProposal, err := s.proposalRepo.GetByID(proposal.ID)
	if err != nil {
		if errors.Is(err, errors.New("proposal not found")) {
			return err
		}
		return fmt.Errorf("failed to get proposal: %w", err)
	}

	proposal.Version = existingProposal.Version + 1
	proposal.UpdatedAt = time.Now()

	existingProposal.Description = proposal.Description
	existingProposal.OrganizationID = proposal.OrganizationID
	existingProposal.Price = proposal.Price
	existingProposal.PublicationDate = proposal.PublicationDate
	existingProposal.UpdatedAt = proposal.UpdatedAt
	existingProposal.Version = proposal.Version
	existingProposal.TenderID = proposal.TenderID
	existingProposal.Status = proposal.Status

	err = s.proposalRepo.Update(existingProposal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("proposal not found")
		}
		return fmt.Errorf("failed to update proposal %w", err)
	}
	return nil
}

func (s *proposalService) Delete(id int64) error {
	return s.proposalRepo.Delete(id)
}

func (s *proposalService) List() ([]*models.Proposal, error) {
	return s.proposalRepo.List()
}

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

func (s *proposalService) Accept(id int64) error {
	proposal, err := s.proposalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("предложение не найдено: %w", err)
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

	if proposal.Status == models.ProposalStatusAccepted {
		return errors.New("предложение уже принято")
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
