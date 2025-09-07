package repository

import (
	"context"
	"database/sql"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
	"github.com/ItsXomyak/scam-list/internal/modules/domains/repository/converters"
	"github.com/ItsXomyak/scam-list/internal/modules/domains/repository/models"
)

type QueriesAdapter struct {
	t *Queries
}

func NewQueriesAdapter(db DBTX) *QueriesAdapter {
	return &QueriesAdapter{
		t: New(db),
	}
}

func (q *QueriesAdapter) DeleteDomain(ctx context.Context, domain string) error {
	return q.t.DeleteDomain(ctx, domain)
}
func (q *QueriesAdapter) GetDomain(ctx context.Context, domain string) (models.DomainDB, error) {
	return q.t.GetDomain(ctx, domain)
}
func (q *QueriesAdapter) GetDomainsByRiskScore(ctx context.Context, arg *entity.GetDomainsByRiskScoreParams) ([]models.DomainDB, error) {

	dbArg := converters.ToDBGetDomainsByRiskScoreParams(arg)
	return q.t.GetDomainsByRiskScore(ctx, dbArg)
}
func (q *QueriesAdapter) GetDomainsByStatus(ctx context.Context, arg *entity.GetDomainsByStatusParams) ([]models.DomainDB, error) {
	dbArgs := converters.ToDBGetDomainsByStatusParams2(arg)
	return q.t.GetDomainsByStatus(ctx, dbArgs)
}
func (q *QueriesAdapter) GetDomainsForRecheck(ctx context.Context, arg *entity.GetDomainsForRecheckParams) ([]models.DomainDB, error) {
	dbArgs := converters.ToDBGetDomainsForRecheckParams(arg)
	return q.t.GetDomainsForRecheck(ctx, dbArgs)
}
func (q *QueriesAdapter) MarkDomainAsScam(ctx context.Context, arg *entity.MarkDomainAsScamParams) (models.DomainDB, error) {
	dbArgs := converters.ToDBMarkDomainAsScamParams(arg)
	return q.t.MarkDomainAsScam(ctx, dbArgs)
}
func (q *QueriesAdapter) UpdateDomainStatus(ctx context.Context, arg *entity.UpdateDomainStatusParams) (models.DomainDB, error) {
	dbArgs := converters.ToDBUpdateDomainStatusParams(arg)
	return q.t.UpdateDomainStatus(ctx, dbArgs)
}
func (q *QueriesAdapter) VerifyDomain(ctx context.Context, arg *entity.VerifyDomainParams) (models.DomainDB, error) {
	dbArgs := converters.ToDBVerifyDomainParams(arg)
	return q.t.VerifyDomain(ctx, dbArgs)
}
func (q *QueriesAdapter) WithTx(tx *sql.Tx) *Queries {
	return q.t.WithTx(tx)
}
func (q *QueriesAdapter) CreateDomain(ctx context.Context, arg *entity.CreateDomainParams) (models.DomainDB, error) {
	dbArgs := converters.ToDBCreateDomainParams(arg)
	return q.t.CreateDomain(ctx, dbArgs)
}