package handler

import (
	"context"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
)

type Verifier interface {
	ProcessDomain(ctx context.Context, url string) (*entity.VerifyDomainResult, error)
}
