package repository

import (
	"database/sql"
	"errors"

	"github.com/itocode21/Tender-Service/internal/models"
)

// OrganizationRepository интерфейс для работы с организациями
type OrganizationRepository interface {
	Create(org *models.Organization) (int64, error)
	GetByID(id int64) (*models.Organization, error)
	Update(org *models.Organization) error
	Delete(id int64) error
	List() ([]*models.Organization, error)
	AddResponsible(orgResp *models.OrganizationResponsible) (int64, error)
	RemoveResponsible(orgResp *models.OrganizationResponsible) error
}

// organizationRepository реализация интерфейса OrganizationRepository
type organizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository создает новый экземпляр organizationRepository
func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// Create создает новую организацию
func (r *organizationRepository) Create(org *models.Organization) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO organizations (name, description, type) VALUES ($1, $2, $3) RETURNING id`,
		org.Name, org.Description, org.Type).Scan(&org.Id)
	return org.Id, err
}

// GetByID получает организацию по ID
func (r *organizationRepository) GetByID(id int64) (*models.Organization, error) {
	row := r.db.QueryRow(`SELECT id, name, description, type, created_at, updated_at FROM organizations WHERE id = $1`, id)
	var org models.Organization
	if err := row.Scan(&org.Id, &org.Name, &org.Description, &org.Type, &org.CreatedAt, &org.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &org, nil
}

// Update обновляет данные организации
func (r *organizationRepository) Update(org *models.Organization) error {
	_, err := r.db.Exec(
		`UPDATE organizations SET name = $1, description = $2, type = $3, updated_at = NOW() WHERE id = $4`,
		org.Name, org.Description, org.Type, org.Id)
	if err != nil {
		return err
	}
	return nil
}

// Delete удаляет организацию по ID
func (r *organizationRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM organizations WHERE id = $1`, id)
	return err
}

// List возвращает список всех организаций
func (r *organizationRepository) List() ([]*models.Organization, error) {
	rows, err := r.db.Query(`SELECT id, name, description, type, created_at, updated_at FROM organizations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []*models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.Id, &org.Name, &org.Description, &org.Type, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		organizations = append(organizations, &org)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return organizations, nil
}

// AddResponsible добавляет ответственного за организацию
func (r *organizationRepository) AddResponsible(orgResp *models.OrganizationResponsible) (int64, error) {
	err := r.db.QueryRow(
		`INSERT INTO organization_responsibles (organization_id, user_id) VALUES ($1, $2) RETURNING id`,
		orgResp.OrganizationID.Id, orgResp.UserID.Id).Scan(&orgResp.Id)
	return orgResp.Id, err
}

// RemoveResponsible удаляет ответственного за организацию
func (r *organizationRepository) RemoveResponsible(orgResp *models.OrganizationResponsible) error {
	_, err := r.db.Exec(`DELETE FROM organization_responsibles WHERE organization_id = $1 AND user_id = $2`,
		orgResp.OrganizationID.Id, orgResp.UserID.Id)
	return err
}
