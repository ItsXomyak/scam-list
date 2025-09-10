package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"

	"github.com/ItsXomyak/scam-list/internal/modules/domains/entity"
	"github.com/ItsXomyak/scam-list/pkg/postgres"
)

type DomainRepository struct {
	pool postgres.PgxPool
}

func NewDomain(pool postgres.PgxPool) *DomainRepository {
	return &DomainRepository{
		pool: pool,
	}
}

func (u *DomainRepository) GetDomain(ctx context.Context, domain string) (*entity.Domain, error) {
	query := `
		SELECT domain, status, company_name, country, scam_sources, scam_type, 
			   verified_by, verification_method, risk_score, reasons, metadata, 
			   created_at, updated_at
		FROM domains 
		WHERE domain = $1`

	d := &entity.Domain{}
	err := u.pool.QueryRow(ctx, query, domain).Scan(
		d.Domain,
		d.Status,
		d.CompanyName,
		d.Country,
		pq.Array(d.ScamSources),
		d.ScamType,
		d.VerifiedBy,
		d.VerificationMethod,
		d.RiskScore,
		pq.Array(d.Reasons),
		pq.Array(d.Metadata),
		d.CreatedAt,
		d.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return d, nil
}

func (u *DomainRepository) CreateDomain(ctx context.Context, arg entity.CreateDomainParams) (*entity.Domain, error) {
	query := `
		INSERT INTO domains (domain, status, company_name, country, scam_sources, 
							scam_type, verified_by, verification_method, risk_score, 
							reasons, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING domain, status, company_name, country, scam_sources, scam_type,
				  verified_by, verification_method, risk_score, reasons, metadata,
				  created_at, updated_at`

	d := &entity.Domain{}
	err := u.pool.QueryRow(ctx, query,
		arg.Domain,
		arg.Status,
		arg.CompanyName,
		arg.Country,
		pq.Array(arg.ScamSources),
		arg.ScamType,
		arg.VerifiedBy,
		arg.VerificationMethod,
		arg.RiskScore,
		pq.Array(arg.Reasons),
		pq.Array(arg.Metadata),
	).Scan(
		d.Domain,
		d.Status,
		d.CompanyName,
		d.Country,
		pq.Array(d.ScamSources),
		d.ScamType,
		d.VerifiedBy,
		d.VerificationMethod,
		d.RiskScore,
		pq.Array(d.Reasons),
		pq.Array(d.Metadata),
		d.CreatedAt,
		d.UpdatedAt,
	)

	return d, err
}
