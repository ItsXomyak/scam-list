package entity

import "encoding/json"

type CreateDomainParams struct {
	Domain             string
	Status             string
	CompanyName        *string
	Country            *string
	ScamSources        []string
	ScamType           *string
	VerifiedBy         *string
	VerificationMethod *string
	RiskScore          *float64
	Reasons            []string
	Metadata           []json.RawMessage
}

// type GetDomainsByRiskScoreParams struct {
// 	RiskScore *string
// 	RiscScore2 *string
// 	Limit      int32
// 	Offset     int32
// }

// type GetDomainsForRecheckParams struct {
// 	LastCheckAt *time.Time
// 	Limit       int32
// }

// type GetDomainsByStatusParams struct {
// 	Status string
// 	Limit  int32
// 	Offset int32
// }

// type MarkDomainAsScamParams struct {
// 	Domain      string
// 	ScamSources []string
// 	ScamType    *string
// 	RiskScore   *string
// 	Reasons     []string
// }

// type UpdateDomainStatusParams struct {
// 	Domain    string
// 	Status    string
// 	RiskScore *string
// 	Reasons   []string
// }

// type VerifyDomainParams struct {
// 	Domain             string
// 	VerifiedAt         *time.Time
// 	VerifiedBy         *string
// 	VerificationMethod *string
// 	ExpiresAt          *time.Time
// 	RiskScore          *string
// 	Reasons            []string
// }
