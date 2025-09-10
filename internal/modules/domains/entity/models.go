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
	VerifiedBy         *string
	VerificationMethod *string
	RiskScore          *string
	Reasons            []string
	Metadata           [][]byte
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

type CheckerResult struct {
	TotalScore float64
}
