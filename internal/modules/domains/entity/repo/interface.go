package repo

import (
	"context"

	. "github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
)

type DomainRepository interface {
	CreateDomain(ctx context.Context, arg CreateDomainParams) (Domain, error)
	DeleteDomain(ctx context.Context, domain string) error
	GetDomain(ctx context.Context, domain string) (Domain, error)
	GetDomainsByRiskScore(ctx context.Context, params GetDomainsByRiskScoreParams) ([]Domain, error)
	GetDomainsForRecheck(ctx context.Context, params GetDomainsForRecheckParams) ([]Domain, error)
	GetDomainsBystatus(ctx context.Context, params GetDomainsByStatusParams) ([]Domain, error)
	UpdateDomainStatus(ctx context.Context, params UpdateDomainStatusParams) (Domain, error)
	MarkDomainAsScam(ctx context.Context, params MarkDomainAsScamParams) (Domain, error)
	VerifyDomain(ctx context.Context, params VerifyDomainParams) (Domain, error)
}