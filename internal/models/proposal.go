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
	TenderID        int64          `json:"tender_id"`
	OrganizationID  int64          `json:"organization_id"`
	Description     string         `json:"description"`
	PublicationDate time.Time      `json:"publication_date"`
	Price           float64        `json:"price"`
	Status          ProposalStatus `json:"status"`
	Version         int            `json:"version"` // Добавили поле Version
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
