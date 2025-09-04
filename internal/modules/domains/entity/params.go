package entity

import "time"

type CreateDomainParams struct {
	Domain             string
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
	LastCheckAt        *time.Time
}

type GetDomainsByRiskScoreParams struct {
	RiskScore *string
	RiscScore2 *string
	Limit      int32
	Offset     int32
}

type GetDomainsForRecheckParams struct {
	LastCheckAt *time.Time
	Limit       int32
}

type GetDomainsByStatusParams struct {
	Status string
	Limit  int32
	Offset int32
}

type MarkDomainAsScamParams struct {
	Domain      string
	ScamSources []string
	ScamType    *string
	RiskScore   *string
	Reasons     []string
}

type UpdateDomainStatusParams struct {
	Domain    string
	Status    string
	RiskScore *string
	Reasons   []string
}

type VerifyDomainParams struct {
	Domain             string
	VerifiedAt         *time.Time
	VerifiedBy         *string
	VerificationMethod *string
	ExpiresAt          *time.Time
	RiskScore          *string
	Reasons            []string
}