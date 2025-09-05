package core

import "github.com/ItsXomyak/scam-list/internal/modules/domains/entity"

func FromDomainToVerifyDomainResult(domain *entity.Domain) *entity.VerifyDomainResult {
	if domain == nil {
		return nil
	}
	return &entity.VerifyDomainResult{
		Domain:             domain.Domain,
		Status:             domain.Status,
		CompanyName:        domain.CompanyName,
		Country:            domain.Country,
		ScamSources:        domain.ScamSources,
		ScamType:           domain.ScamType,
		VerifiedAt:         domain.VerifiedAt,
		VerifiedBy:         domain.VerifiedBy,
		VerificationMethod: domain.VerificationMethod,
		ExpiresAt:          domain.ExpiresAt,
		RiskScore:          domain.RiskScore,
		Reasons:            domain.Reasons,
		Metadata:           domain.Metadata,
		CreatedAt:          domain.CreatedAt,
		UpdatedAt:          domain.UpdatedAt,
		LastCheckAt:        domain.LastCheckAt,
	}
}