package repository

import (
	"database/sql"

	"github.com/sqlc-dev/pqtype"
)

type СreateDomainParams struct {
	Domain             string
	CompanyName        sql.NullString
	Country            sql.NullString
	ScamSources        []string
	ScamType           sql.NullString
	VerifiedAt         sql.NullTime
	VerifiedBy         sql.NullString
	VerificationMethod sql.NullString
	ExpiresAt          sql.NullTime
	RiskScore          sql.NullString
	Reasons            []string
	Metadata           pqtype.NullRawMessage
	LastCheckAt        sql.NullTime
}

type GetDomainsByRiskScoreParams struct {
	RiskScore   sql.NullString
	RiskScore_2 sql.NullString
	Limit       int32
	Offset      int32
}

type GetDomainsByStatusParams struct {
	Status string
	Limit  int32
	Offset int32
}

type GetDomainsForRecheckParams struct {
	LastCheckAt sql.NullTime
	Limit       int32
}

type MarkDomainAsScamParams struct {
	Domain      string
	ScamSources []string
	ScamType    sql.NullString
	RiskScore   sql.NullString
	Reasons     []string
}

type UpdateDomainStatusParams struct {
	Domain    string
	Status    string
	RiskScore sql.NullString
	Reasons   []string
}

type VerifyDomainParams struct {
	Domain             string
	VerifiedAt         sql.NullTime
	VerifiedBy         sql.NullString
	VerificationMethod sql.NullString
	ExpiresAt          sql.NullTime
	RiskScore          sql.NullString
	Reasons            []string
}

const СreateDomain = `-- name: СreateDomain :one
INSERT INTO domains (
    domain, company_name, country, scam_sources, scam_type,
    verified_at, verified_by, verification_method, expires_at,
    risk_score, reasons, metadata, last_check_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at
`

const DeleteDomain = `-- name: DeleteDomain :exec
DELETE FROM domains WHERE domain = $1
`

const GetDomain = `-- name: GetDomain :one
SELECT domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at FROM domains WHERE domain = $1
`

const GetDomainsByRiskScore = `-- name: GetDomainsByRiskScore :many
SELECT domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at FROM domains 
WHERE risk_score BETWEEN $1 AND $2 
ORDER BY risk_score DESC 
LIMIT $3 OFFSET $4
`

const GetDomainsByStatus = `-- name: GetDomainsByStatus :many
SELECT domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at FROM domains 
WHERE status = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3
`

const GetDomainsForRecheck = `-- name: GetDomainsForRecheck :many
SELECT domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at FROM domains 
WHERE last_check_at IS NULL 
OR last_check_at < $1 
ORDER BY last_check_at ASC NULLS FIRST
LIMIT $2
`

const MarkDomainAsScam = `-- name: MarkDomainAsScam :one
UPDATE domains 
SET status = 'scam',
    scam_sources = $2,
    scam_type = $3,
    risk_score = $4,
    reasons = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at
`

const UpdateDomainStatus = `-- name: UpdateDomainStatus :one
UPDATE domains 
SET status = $2, risk_score = $3, reasons = $4, updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at
`

const VerifyDomain = `-- name: VerifyDomain :one
UPDATE domains 
SET status = 'verified', 
    verified_at = COALESCE($2, CURRENT_TIMESTAMP),
    verified_by = $3,
    verification_method = $4,
    expires_at = $5,
    risk_score = $6,
    reasons = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING domain, status, company_name, country, scam_sources, scam_type, verified_at, verified_by, verification_method, expires_at, risk_score, reasons, metadata, created_at, updated_at, last_check_at
`