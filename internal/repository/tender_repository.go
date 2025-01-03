package repository

import (
	"database/sql"
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
)

type TenderRepository interface {
	Create(tender *models.Tender) (int64, error)
	GetByID(id int64) (*models.Tender, error)
	Update(tender *models.Tender) error
	Delete(id int64) error
	List() ([]*models.Tender, error)
	ListByOrganizationID(organizationId int64) ([]*models.Tender, error)
}

type tenderRepository struct {
	db *sql.DB
}

func NewTenderRepository(db *sql.DB) TenderRepository {
	return &tenderRepository{db: db}
}

func (r *tenderRepository) Create(tender *models.Tender) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO tenders (name, description, organization_id, publication_date, end_date, status, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		tender.Name, tender.Description, tender.OrganizationID, tender.PublicationDate, tender.EndDate, models.TenderStatusCreated, tender.Version).Scan(&tender.ID)
	if err != nil {
		return 0, err
	}
	tender.Status = models.TenderStatusCreated
	return tender.ID, nil
}

func (r *tenderRepository) GetByID(id int64) (*models.Tender, error) {
	row := r.db.QueryRow(`SELECT id, name, description, organization_id, publication_date, end_date, status, version, created_at, updated_at FROM tenders WHERE id = $1`, id)
	var tender models.Tender
	if err := row.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.OrganizationID, &tender.PublicationDate, &tender.EndDate, &tender.Status, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}
	return &tender, nil
}

func (r *tenderRepository) Update(tender *models.Tender) error {
	_, err := r.db.Exec(
		`UPDATE tenders SET name = $1, description = $2, organization_id = $3, publication_date = $4, end_date = $5, status = $6, version = $7, updated_at = NOW() WHERE id = $8`,
		tender.Name, tender.Description, tender.OrganizationID, tender.PublicationDate, tender.EndDate, tender.Status, tender.Version, tender.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *tenderRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tenders WHERE id = $1`, id)
	return err
}

func (r *tenderRepository) List() ([]*models.Tender, error) {
	rows, err := r.db.Query(`SELECT id, name, description, organization_id, publication_date, end_date, status, version, created_at, updated_at FROM tenders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []*models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.OrganizationID, &tender.PublicationDate, &tender.EndDate, &tender.Status, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tenders, nil
}

func (r *tenderRepository) ListByOrganizationID(organizationId int64) ([]*models.Tender, error) {
	rows, err := r.db.Query(`SELECT id, name, description, organization_id, publication_date, end_date, status, version, created_at, updated_at FROM tenders WHERE organization_id = $1`, organizationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []*models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.OrganizationID, &tender.PublicationDate, &tender.EndDate, &tender.Status, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tenders, nil
}
