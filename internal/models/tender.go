package models

import "time"

// TenderStatus представляет собой статус тендера.
type TenderStatus string

const (
	TenderStatusDraft     TenderStatus = "draft"
	TenderStatusPublished TenderStatus = "published"
	TenderStatusClosed    TenderStatus = "closed"
	TenderStatusCancelled TenderStatus = "cancelled"
)

type Tender struct {
	ID              int64        `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	OrganizationID  int64        `json:"organization_id"`  // ID организации, которая создала тендер
	PublicationDate time.Time    `json:"publication_date"` // Дата публикации тендера
	EndDate         time.Time    `json:"end_date"`         // Дата окончания приема заявок
	Status          TenderStatus `json:"status"`           // Статус тендера
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}
