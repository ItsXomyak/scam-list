package repository

import (
	"context"
	"database/sql"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
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
func (q *QueriesAdapter) GetDomain(ctx context.Context, domain string) (DomainDB, error) {
	return q.t.GetDomain(ctx, domain)
}
func (q *QueriesAdapter) GetDomainsByRiskScore(ctx context.Context, arg *entity.GetDomainsByRiskScoreParams) ([]DomainDB, error) {

	dbArg := ToDBGetDomainsByRiskScoreParams(arg)
	return q.t.GetDomainsByRiskScore(ctx, dbArg)
}
func (q *QueriesAdapter) GetDomainsByStatus(ctx context.Context, arg *entity.GetDomainsByStatusParams) ([]DomainDB, error) {
	dbArgs := ToDBGetDomainsByStatusParams2(arg)
	return q.t.GetDomainsByStatus(ctx, dbArgs)
}
func (q *QueriesAdapter) GetDomainsForRecheck(ctx context.Context, arg *entity.GetDomainsForRecheckParams) ([]DomainDB, error) {
	dbArgs := ToDBGetDomainsForRecheckParams(arg)
	return q.t.GetDomainsForRecheck(ctx, dbArgs)
}
func (q *QueriesAdapter) MarkDomainAsScam(ctx context.Context, arg *entity.MarkDomainAsScamParams) (DomainDB, error) {
	dbArgs := ToDBMarkDomainAsScamParams(arg)
	return q.t.MarkDomainAsScam(ctx, dbArgs)
}
func (q *QueriesAdapter) UpdateDomainStatus(ctx context.Context, arg UpdateDomainStatusParams) (DomainDB, error) {
	return q.t.UpdateDomainStatus(ctx, arg)
}
func (q *QueriesAdapter) VerifyDomain(ctx context.Context, arg *entity.VerifyDomainParams) (DomainDB, error) {
	dbArgs := ToDBVerifyDomainParams(arg)
	return q.t.VerifyDomain(ctx, dbArgs)
}
func (q *QueriesAdapter) WithTx(tx *sql.Tx) *Queries {
	return q.t.WithTx(tx)
}
func (q *QueriesAdapter) CreateDomain(ctx context.Context, arg *entity.CreateDomainParams) (DomainDB, error) {
	dbArgs := ToDBCreateDomainParams(arg)
	return q.t.CreateDomain(ctx, dbArgs)
}