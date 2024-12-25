package repository

import (
	"database/sql"
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
)

// TenderRepository интерфейс для работы с тендерами
type TenderRepository interface {
	Create(tender *models.Tender) (int64, error)
	GetByID(id int64) (*models.Tender, error)
	GetByOrganizationID(organizationID int64) ([]*models.Tender, error)
	Update(tender *models.Tender) error
	Delete(id int64) error
	List() ([]*models.Tender, error)
}

// tenderRepository реализация интерфейса TenderRepository
type tenderRepository struct {
	db *sql.DB
}

// NewTenderRepository создает новый экземпляр tenderRepository
func NewTenderRepository(db *sql.DB) TenderRepository {
	return &tenderRepository{db: db}
}

// Create создает новый тендер
func (r *tenderRepository) Create(tender *models.Tender) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO tenders (name, description, organization_id, publication_date, end_date, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		tender.Name, tender.Description, tender.OrganizationID, tender.PublicationDate, tender.EndDate, tender.Status).Scan(&tender.ID)
	return tender.ID, err
}

// GetByID получает тендер по ID
func (r *tenderRepository) GetByID(id int64) (*models.Tender, error) {
	row := r.db.QueryRow(`SELECT id, name, description, organization_id, publication_date, end_date, status, created_at, updated_at FROM tenders WHERE id = $1`, id)
	var tender models.Tender
	if err := row.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.OrganizationID, &tender.PublicationDate, &tender.EndDate, &tender.Status, &tender.CreatedAt, &tender.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}
	return &tender, nil
}

// GetByOrganizationID получает все тендеры по id организации
func (r *tenderRepository) GetByOrganizationID(organizationID int64) ([]*models.Tender, error) {
	rows, err := r.db.Query(`SELECT id, name, description, organization_id, publication_date, end_date, status, created_at, updated_at FROM tenders WHERE organization_id = $1`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []*models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.OrganizationID, &tender.PublicationDate, &tender.EndDate, &tender.Status, &tender.CreatedAt, &tender.UpdatedAt); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tenders, nil
}

// Update обновляет данные тендера
func (r *tenderRepository) Update(tender *models.Tender) error {
	_, err := r.db.Exec(
		`UPDATE tenders SET name = $1, description = $2, organization_id = $3, publication_date = $4, end_date = $5, status = $6, updated_at = NOW() WHERE id = $7`,
		tender.Name, tender.Description, tender.OrganizationID, tender.PublicationDate, tender.EndDate, tender.Status, tender.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete удаляет тендер по ID
func (r *tenderRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tenders WHERE id = $1`, id)
	return err
}

// List возвращает список всех тендеров
func (r *tenderRepository) List() ([]*models.Tender, error) {
	rows, err := r.db.Query(`SELECT id, name, description, organization_id, publication_date, end_date, status, created_at, updated_at FROM tenders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []*models.Tender
	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.OrganizationID, &tender.PublicationDate, &tender.EndDate, &tender.Status, &tender.CreatedAt, &tender.UpdatedAt); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tenders, nil
}
