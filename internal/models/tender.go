package models

import "time"

// TenderStatus представляет собой статус тендера.
type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "created"
	TenderStatusPublished TenderStatus = "published"
	TenderStatusClosed    TenderStatus = "closed"
	TenderStatusCancelled TenderStatus = "cancelled"
)

type Tender struct {
	ID              int64        `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	OrganizationID  int64        `json:"organization_id"`
	PublicationDate time.Time    `json:"publication_date"`
	EndDate         time.Time    `json:"end_date"`
	Status          TenderStatus `json:"status"`
	Version         int          `json:"version"` // Добавили поле Version
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}
