package dto

import (
	"encoding/json"
	"time"

	"github.com/ItsXomyak/scam-list/internal/domain/entity"
)

type CreateDomainRequest struct {
	Domain             string            `json:"domain"`
	Status             string            `json:"status"`
	CompanyName        *string           `json:"company_name,omitempty"`
	Country            *string           `json:"country,omitempty"`
	ScamSources        []string          `json:"scam_sources,omitempty"`
	ScamType           *string           `json:"scam_type,omitempty"`
	VerifiedBy         *string           `json:"verified_by,omitempty"`
	VerificationMethod *string           `json:"verification_method,omitempty"`
	RiskScore          *float64          `json:"risk_score,omitempty"`
	Reasons            []string          `json:"reasons,omitempty"`
	Metadata           []json.RawMessage `json:"metadata,omitempty"`
}

type DomainResponse struct {
	Domain             string            `json:"domain"`
	Status             string            `json:"status"`
	CompanyName        *string           `json:"company_name"`
	Country            *string           `json:"country"`
	ScamSources        []string          `json:"scam_sources"`
	ScamType           *string           `json:"scam_type"`
	VerifiedBy         *string           `json:"verified_by"`
	VerificationMethod *string           `json:"verification_method"`
	RiskScore          *float64          `json:"risk_score"`
	Reasons            []string          `json:"reasons"`
	Metadata           []json.RawMessage `json:"metadata"`
	CreatedAt          *string           `json:"created_at"`
	UpdatedAt          *string           `json:"updated_at"`
}

type UpdateDomainRequest struct {
	Status             *string           `json:"status,omitempty"`
	CompanyName        *string           `json:"company_name,omitempty"`
	Country            *string           `json:"country,omitempty"`
	ScamSources        []string          `json:"scam_sources,omitempty"`
	ScamType           *string           `json:"scam_type,omitempty"`
	VerifiedBy         *string           `json:"verified_by,omitempty"`
	VerificationMethod *string           `json:"verification_method,omitempty"`
	RiskScore          *float64          `json:"risk_score,omitempty"`
	Reasons            []string          `json:"reasons,omitempty"`
	Metadata           []json.RawMessage `json:"metadata,omitempty"` // или *[]byte[] если оставляешь [][]byte
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

func ToBatchDomainResponse(d []*entity.Domain) []*DomainResponse {
	res := make([]*DomainResponse, len(d))

	for _, v := range d {
		res = append(res, ToDomainResponse(v))
	}

	return res
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
