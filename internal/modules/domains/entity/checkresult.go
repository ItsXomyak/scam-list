package entity

import "time"

// VerifyDomainResult represents the result of verifying a domain that we returns the user.
type VerifyDomainResult struct {
	Domain        string          `json:"domain"`
	Status        string          `json:"status"`
	ScamType      string          `json:"scam_type"`
	RiskScore     float64         `json:"risk_score"`
	CompanyName   string          `json:"company_name"`
	Country       string          `json:"country"`
	VerifiedBy    string          `json:"verified_by"`
	VerifiedAt    time.Time       `json:"verified_at"`
	ModuleResults []*ModuleResult `json:"module_results"`
}

type ModuleResult struct {
	ModuleName  string         `json:"module_name"`
	Description string         `json:"description"`
	RiskScore   float64        `json:"risk_score"`
	Metadata    map[string]any `json:"metadata"`
}
