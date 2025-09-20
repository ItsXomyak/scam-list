package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/ItsXomyak/scam-list/internal/domain/entity"
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

// =====================
// helpers for metadata
// =====================

// packMetadata: [][]byte (каждый элемент — валидный JSON) -> JSON-массив []byte
func packMetadata(src [][]byte) ([]byte, error) {
	if src == nil {
		return nil, nil // позволим NULL в БД
	}
	arr := make([]json.RawMessage, len(src))
	for i := range src {
		if !json.Valid(src[i]) {
			return nil, fmt.Errorf("metadata[%d] is not valid JSON", i)
		}
		arr[i] = json.RawMessage(src[i])
	}
	return json.Marshal(arr)
}

// unpackMetadata: JSON-массив []byte -> [][]byte
func unpackMetadata(b []byte) ([][]byte, error) {
	if len(b) == 0 {
		return nil, nil
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(b, &arr); err != nil {
		// если вдруг в БД не массив, вернём ошибку явную
		return nil, fmt.Errorf("metadata is not a JSON array: %w", err)
	}
	out := make([][]byte, len(arr))
	for i := range arr {
		// делаем копию, чтобы не держать ссылки на один буфер
		if arr[i] == nil {
			out[i] = nil
			continue
		}
		cp := make([]byte, len(arr[i]))
		copy(cp, arr[i])
		out[i] = cp
	}
	return out, nil
}

// =====================
// CRUD
// =====================

func (u *DomainRepository) CreateDomain(ctx context.Context, arg entity.CreateDomainParams) (*entity.Domain, error) {
	// подготовим metadata
	mdJSON, err := packMetadata(arg.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO domains (
			domain, status, company_name, country, scam_sources,
			scam_type, verified_by, verification_method, risk_score,
			reasons, metadata
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9::numeric,
			$10, $11::jsonb
		)
		RETURNING
			domain,
			status,
			company_name,
			country,
			scam_sources,
			scam_type,
			verified_by,
			verification_method,
			risk_score::text,
			reasons,
			metadata,
			created_at,
			updated_at
	`

	var (
		res                    entity.Domain
		metadataRaw            []byte
		riskScoreText, company *string
		country                *string
		scamType               *string
		verifiedBy             *string
		verificationMethod     *string
		createdAt, updatedAt   *time.Time
		scamSources, reasons   []string
	)

	err = u.pool.QueryRow(ctx, query,
		arg.Domain,
		arg.Status,
		arg.CompanyName,        // nullable
		arg.Country,            // nullable
		arg.ScamSources,        // text[]
		arg.ScamType,           // nullable
		arg.VerifiedBy,         // nullable
		arg.VerificationMethod, // nullable
		arg.RiskScore,          // ::numeric (nullable строка)
		arg.Reasons,            // text[]
		mdJSON,                 // ::jsonb
	).Scan(
		&res.Domain,
		&res.Status,
		&company,
		&country,
		&scamSources,
		&scamType,
		&verifiedBy,
		&verificationMethod,
		&riskScoreText,
		&reasons,
		&metadataRaw,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	md, err := unpackMetadata(metadataRaw)
	if err != nil {
		return nil, err
	}

	res.CompanyName = company
	res.Country = country
	res.ScamSources = scamSources
	res.ScamType = scamType
	res.VerifiedBy = verifiedBy
	res.VerificationMethod = verificationMethod
	res.RiskScore = riskScoreText
	res.Reasons = reasons
	res.Metadata = md
	res.CreatedAt = createdAt
	res.UpdatedAt = updatedAt

	return &res, nil
}

func (u *DomainRepository) GetDomain(ctx context.Context, domain string) (*entity.Domain, error) {
	query := `
		SELECT
			domain,
			status,
			company_name,
			country,
			scam_sources,
			scam_type,
			verified_by,
			verification_method,
			risk_score::text,
			reasons,
			metadata,
			created_at,
			updated_at
		FROM domains
		WHERE domain = $1
	`

	var (
		res                    entity.Domain
		metadataRaw            []byte
		riskScoreText, company *string
		country                *string
		scamType               *string
		verifiedBy             *string
		verificationMethod     *string
		createdAt, updatedAt   *time.Time
		scamSources, reasons   []string
	)

	err := u.pool.QueryRow(ctx, query, domain).Scan(
		&res.Domain,
		&res.Status,
		&company,
		&country,
		&scamSources,
		&scamType,
		&verifiedBy,
		&verificationMethod,
		&riskScoreText,
		&reasons,
		&metadataRaw,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	md, err := unpackMetadata(metadataRaw)
	if err != nil {
		return nil, err
	}

	res.CompanyName = company
	res.Country = country
	res.ScamSources = scamSources
	res.ScamType = scamType
	res.VerifiedBy = verifiedBy
	res.VerificationMethod = verificationMethod
	res.RiskScore = riskScoreText
	res.Reasons = reasons
	res.Metadata = md
	res.CreatedAt = createdAt
	res.UpdatedAt = updatedAt

	return &res, nil
}

func (u *DomainRepository) GetAllDomains(ctx context.Context) ([]*entity.Domain, error) {
	query := `
		SELECT
			domain,
			status,
			company_name,
			country,
			scam_sources,
			scam_type,
			verified_by,
			verification_method,
			risk_score::text,
			reasons,
			metadata,
			created_at,
			updated_at
		FROM domains
		ORDER BY created_at DESC
	`

	rows, err := u.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Domain

	for rows.Next() {
		var (
			d                      entity.Domain
			metadataRaw            []byte
			riskScoreText, company *string
			country                *string
			scamType               *string
			verifiedBy             *string
			verificationMethod     *string
			createdAt, updatedAt   *time.Time
			scamSources, reasons   []string
		)

		if err := rows.Scan(
			&d.Domain,
			&d.Status,
			&company,
			&country,
			&scamSources,
			&scamType,
			&verifiedBy,
			&verificationMethod,
			&riskScoreText,
			&reasons,
			&metadataRaw,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		md, err := unpackMetadata(metadataRaw)
		if err != nil {
			return nil, err
		}

		d.CompanyName = company
		d.Country = country
		d.ScamSources = scamSources
		d.ScamType = scamType
		d.VerifiedBy = verifiedBy
		d.VerificationMethod = verificationMethod
		d.RiskScore = riskScoreText
		d.Reasons = reasons
		d.Metadata = md
		d.CreatedAt = createdAt
		d.UpdatedAt = updatedAt

		out = append(out, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (u *DomainRepository) UpdateDomain(ctx context.Context, updated *entity.Domain) (*entity.Domain, error) {
	mdJSON, err := packMetadata(updated.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE domains SET
			status = $2,
			company_name = $3,
			country = $4,
			scam_sources = $5,
			scam_type = $6,
			verified_by = $7,
			verification_method = $8,
			risk_score = $9::numeric,
			reasons = $10,
			metadata = $11,
			updated_at = NOW()
		WHERE domain = $1
		RETURNING
			domain,
			status,
			company_name,
			country,
			scam_sources,
			scam_type,
			verified_by,
			verification_method,
			risk_score::text,
			reasons,
			metadata,
			created_at,
			updated_at
	`

	var (
		res                    entity.Domain
		metadataRaw            []byte
		riskScoreText, company *string
		country                *string
		scamType               *string
		verifiedBy             *string
		verificationMethod     *string
		createdAt, updatedAt   *time.Time
		scamSources, reasons   []string
	)

	err = u.pool.QueryRow(ctx, query,
		updated.Domain,
		updated.Status,
		updated.CompanyName,
		updated.Country,
		updated.ScamSources,
		updated.ScamType,
		updated.VerifiedBy,
		updated.VerificationMethod,
		updated.RiskScore, // ::numeric
		updated.Reasons,
		mdJSON,
	).Scan(
		&res.Domain,
		&res.Status,
		&company,
		&country,
		&scamSources,
		&scamType,
		&verifiedBy,
		&verificationMethod,
		&riskScoreText,
		&reasons,
		&metadataRaw,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	md, err := unpackMetadata(metadataRaw)
	if err != nil {
		return nil, err
	}

	res.CompanyName = company
	res.Country = country
	res.ScamSources = scamSources
	res.ScamType = scamType
	res.VerifiedBy = verifiedBy
	res.VerificationMethod = verificationMethod
	res.RiskScore = riskScoreText
	res.Reasons = reasons
	res.Metadata = md
	res.CreatedAt = createdAt
	res.UpdatedAt = updatedAt

	return &res, nil
}

func (u *DomainRepository) DeleteDomain(ctx context.Context, domain string) error {
	cmd, err := u.pool.Exec(ctx, `DELETE FROM domains WHERE domain = $1`, domain)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
