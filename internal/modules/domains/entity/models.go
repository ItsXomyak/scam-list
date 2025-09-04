package entity

import (
	"time"
)

type Domain struct {
	Domain             string
	Status             string
	CompanyName        *string
	Country            *string
	ScamSources        []string
	ScamType           *string
	VerifiedAt         *time.Time
	VerifiedBy         *string
	VerificationMethod *string
	ExpiresAt          *time.Time
	RiskScore          *string
	Reasons            []string
	Metadata           [][]byte
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
	LastCheckAt        *time.Time
}

type PendingModeration struct {
	Domain         string
	CheckID        string
	Reasons        []string
	SourceModules  []string
	Priority       *int32
	Status         *string
	AssignedTo     *string
	SubmittedAt    *time.Time
	ResolvedAt     *time.Time
	ModeratorNotes *string
	CreatedAt      *time.Time
}

