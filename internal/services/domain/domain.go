package domain

import (
	"context"

	"github.com/ItsXomyak/scam-list/internal/domain/entity"
)

type DomainRepository interface {
	CreateDomain(ctx context.Context, arg entity.CreateDomainParams) (*entity.Domain, error)
	GetDomain(ctx context.Context, domain string) (*entity.Domain, error)
}

type DomainService struct {
	repo DomainRepository
}

func NewDomainService(repo DomainRepository) *DomainService {
	return &DomainService{repo: repo}
}

func (s *DomainService) CreateDomain(ctx context.Context, params entity.CreateDomainParams) (*entity.Domain, error) {
	return s.repo.CreateDomain(ctx, params)
}

func (s *DomainService) GetDomain(ctx context.Context, domain string) (*entity.Domain, error) {
	return s.repo.GetDomain(ctx, domain)
}
