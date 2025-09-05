package repo

import (
	"context"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
)

type DomainRepository interface {
	CreateDomain(ctx context.Context, arg entity.CreateDomainParams) (entity.Domain, error)
	DeleteDomain(ctx context.Context, domain string) error
	GetDomain(ctx context.Context, domain string) (entity.Domain, error)
	GetDomainsByRiskScore(ctx context.Context, params entity.GetDomainsByRiskScoreParams) ([]entity.Domain, error)
	GetDomainsForRecheck(ctx context.Context, params entity.GetDomainsForRecheckParams) ([]entity.Domain, error)
	GetDomainsBystatus(ctx context.Context, params entity.GetDomainsByStatusParams) ([]entity.Domain, error)
	UpdateDomainStatus(ctx context.Context, params entity.UpdateDomainStatusParams) (entity.Domain, error)
	MarkDomainAsScam(ctx context.Context, params entity.MarkDomainAsScamParams) (entity.Domain, error)
	VerifyDomain(ctx context.Context, params entity.VerifyDomainParams) (entity.Domain, error)
}