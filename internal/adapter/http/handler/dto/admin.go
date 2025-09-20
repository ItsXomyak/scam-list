package dto

import (
	"time"

	"github.com/ItsXomyak/scam-list/internal/domain/entity"
)

type CreateDomainRequest struct {
	Domain             string   `json:"domain"`
	Status             string   `json:"status"`
	CompanyName        *string  `json:"company_name,omitempty"`
	Country            *string  `json:"country,omitempty"`
	ScamSources        []string `json:"scam_sources,omitempty"`
	ScamType           *string  `json:"scam_type,omitempty"`
	VerifiedBy         *string  `json:"verified_by,omitempty"`
	VerificationMethod *string  `json:"verification_method,omitempty"`
	RiskScore          *string  `json:"risk_score,omitempty"`
	Reasons            []string `json:"reasons,omitempty"`
	Metadata           [][]byte `json:"metadata,omitempty"`
}

type DomainResponse struct {
	Domain             string   `json:"domain"`
	Status             string   `json:"status"`
	CompanyName        *string  `json:"company_name"`
	Country            *string  `json:"country"`
	ScamSources        []string `json:"scam_sources"`
	ScamType           *string  `json:"scam_type"`
	VerifiedBy         *string  `json:"verified_by"`
	VerificationMethod *string  `json:"verification_method"`
	RiskScore          *string  `json:"risk_score"`
	Reasons            []string `json:"reasons"`
	Metadata           [][]byte `json:"metadata"`
	CreatedAt          *string  `json:"created_at"`
	UpdatedAt          *string  `json:"updated_at"`
}

func FromCreateRequestToInternal(req *CreateDomainRequest) *entity.CreateDomainParams {
	if req == nil {
		return nil
	}
	return &entity.CreateDomainParams{
		Domain:             req.Domain,
		Status:             req.Status,
		CompanyName:        req.CompanyName,
		Country:            req.Country,
		ScamSources:        req.ScamSources,
		ScamType:           req.ScamType,
		VerifiedBy:         req.VerifiedBy,
		VerificationMethod: req.VerificationMethod,
		RiskScore:          req.RiskScore,
		Reasons:            req.Reasons,
		Metadata:           req.Metadata,
	}
}

func ToDomainResponse(d *entity.Domain) *DomainResponse {
	if d == nil {
		return nil
	}

	var createdAtStr, updatedAtStr *string
	if d.CreatedAt != nil {
		s := d.CreatedAt.Format(time.RFC3339)
		createdAtStr = &s
	}
	if d.UpdatedAt != nil {
		s := d.UpdatedAt.Format(time.RFC3339)
		updatedAtStr = &s
	}

	return &DomainResponse{
		Domain:             d.Domain,
		Status:             d.Status,
		CompanyName:        d.CompanyName,
		Country:            d.Country,
		ScamSources:        d.ScamSources,
		ScamType:           d.ScamType,
		VerifiedBy:         d.VerifiedBy,
		VerificationMethod: d.VerificationMethod,
		RiskScore:          d.RiskScore,
		Reasons:            d.Reasons,
		Metadata:           d.Metadata,
		CreatedAt:          createdAtStr,
		UpdatedAt:          updatedAtStr,
	}
}
