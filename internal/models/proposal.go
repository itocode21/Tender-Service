package models

import "time"

// ProposalStatus представляет собой статус предложения.
type ProposalStatus string

const (
	ProposalStatusDraft     ProposalStatus = "draft"
	ProposalStatusPublished ProposalStatus = "published"
	ProposalStatusAccepted  ProposalStatus = "accepted"
	ProposalStatusRejected  ProposalStatus = "rejected"
	ProposalStatusCancelled ProposalStatus = "cancelled"
)

type Proposal struct {
	ID              int64          `json:"id"`
	TenderID        int64          `json:"tender_id"`       // ID тендера, к которому относится предложение
	OrganizationID  int64          `json:"organization_id"` // ID организации, которая делает предложение
	Description     string         `json:"description"`
	PublicationDate time.Time      `json:"publication_date"` // Дата публикации предложения
	Price           float64        `json:"price"`            // Цена предложения
	Status          ProposalStatus `json:"status"`           // Статус предложения
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
