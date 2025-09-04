package service

import (
	"github.com/ItsXomyak/scam-list/internal/modules/domains/repository"
)

type DomainService struct {
	repo *repository.Queries
}


func NewDomainService(repo *repository.Queries) *DomainService {
	return &DomainService{repo: repo}
}

