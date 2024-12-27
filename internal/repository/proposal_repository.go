package repository

import (
	"database/sql"
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
)

// ProposalRepository интерфейс для работы с предложениями
type ProposalRepository interface {
	Create(proposal *models.Proposal) (int64, error)
	GetByID(id int64) (*models.Proposal, error)
	GetByTenderID(tenderID int64) ([]*models.Proposal, error)
	Update(proposal *models.Proposal) error
	Delete(id int64) error
	List() ([]*models.Proposal, error)
}

// proposalRepository реализация интерфейса ProposalRepository
type proposalRepository struct {
	db *sql.DB
}

// NewProposalRepository создает новый экземпляр proposalRepository
func NewProposalRepository(db *sql.DB) ProposalRepository {
	return &proposalRepository{db: db}
}

// Create создает новое предложение
func (r *proposalRepository) Create(proposal *models.Proposal) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO proposals (tender_id, organization_id, description, publication_date, price, status, version)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		proposal.TenderID, proposal.OrganizationID, proposal.Description, proposal.PublicationDate, proposal.Price, string(models.ProposalStatusCreated), proposal.Version).Scan(&proposal.ID)
	if err != nil {
		return 0, err
	}
	proposal.Status = models.ProposalStatusCreated
	return proposal.ID, nil
}

// GetByID получает предложение по ID
func (r *proposalRepository) GetByID(id int64) (*models.Proposal, error) {
	row := r.db.QueryRow(
		`SELECT id, tender_id, organization_id, description, publication_date, price, status, version, created_at, updated_at
        FROM proposals WHERE id = $1`, id)
	var proposal models.Proposal
	if err := row.Scan(&proposal.ID, &proposal.TenderID, &proposal.OrganizationID, &proposal.Description, &proposal.PublicationDate, &proposal.Price, &proposal.Status, &proposal.Version, &proposal.CreatedAt, &proposal.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("proposal not found")
		}
		return nil, err
	}
	return &proposal, nil
}

// GetByTenderID получает все предложения для определенного тендера
func (r *proposalRepository) GetByTenderID(tenderID int64) ([]*models.Proposal, error) {
	rows, err := r.db.Query(
		`SELECT id, tender_id, organization_id, description, publication_date, price, status, version, created_at, updated_at
        FROM proposals WHERE tender_id = $1`, tenderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proposals []*models.Proposal
	for rows.Next() {
		var proposal models.Proposal
		if err := rows.Scan(&proposal.ID, &proposal.TenderID, &proposal.OrganizationID, &proposal.Description, &proposal.PublicationDate, &proposal.Price, &proposal.Status, &proposal.Version, &proposal.CreatedAt, &proposal.UpdatedAt); err != nil {
			return nil, err
		}
		proposals = append(proposals, &proposal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return proposals, nil
}

// Update обновляет данные предложения
func (r *proposalRepository) Update(proposal *models.Proposal) error {
	_, err := r.db.Exec(
		`UPDATE proposals SET tender_id = $1, organization_id = $2, description = $3, publication_date = $4, price = $5, status = $6, version = $7, updated_at = NOW() WHERE id = $8`,
		proposal.TenderID, proposal.OrganizationID, proposal.Description, proposal.PublicationDate, proposal.Price, proposal.Status, proposal.Version, proposal.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete удаляет предложение по ID
func (r *proposalRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM proposals WHERE id = $1`, id)
	return err
}

// List возвращает список всех предложений
func (r *proposalRepository) List() ([]*models.Proposal, error) {
	rows, err := r.db.Query(
		`SELECT id, tender_id, organization_id, description, publication_date, price, status, version, created_at, updated_at
        FROM proposals`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proposals []*models.Proposal
	for rows.Next() {
		var proposal models.Proposal
		if err := rows.Scan(&proposal.ID, &proposal.TenderID, &proposal.OrganizationID, &proposal.Description, &proposal.PublicationDate, &proposal.Price, &proposal.Status, &proposal.Version, &proposal.CreatedAt, &proposal.UpdatedAt); err != nil {
			return nil, err
		}
		proposals = append(proposals, &proposal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return proposals, nil
}
