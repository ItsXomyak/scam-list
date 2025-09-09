-- name: CreateDomain :one
INSERT INTO domains (
    domain, company_name, country, scam_sources, scam_type,
    verified_at, verified_by, verification_method, expires_at,
    risk_score, reasons, metadata, last_check_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetDomain :one
SELECT * FROM domains WHERE domain = $1;

-- name: GetDomainsByStatus :many
SELECT * FROM domains 
WHERE status = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetDomainsByRiskScore :many
SELECT * FROM domains 
WHERE risk_score BETWEEN $1 AND $2 
ORDER BY risk_score DESC 
LIMIT $3 OFFSET $4;

-- name: UpdateDomainStatus :one
UPDATE domains 
SET status = $2, risk_score = $3, reasons = $4, updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING *;

-- name: VerifyDomain :one
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
RETURNING *;

-- name: MarkDomainAsScam :one
UPDATE domains 
SET status = 'scam',
    scam_sources = $2,
    scam_type = $3,
    risk_score = $4,
    reasons = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING *;

-- name: DeleteDomain :exec
DELETE FROM domains WHERE domain = $1;

-- name: GetDomainsForRecheck :many
SELECT * FROM domains 
WHERE last_check_at IS NULL 
OR last_check_at < $1 
ORDER BY last_check_at ASC NULLS FIRST
LIMIT $2;


