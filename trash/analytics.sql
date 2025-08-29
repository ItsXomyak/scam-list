-- name: GetDomainStatsByCountry :many
SELECT 
    country,
    COUNT(*) as total_count,
    COUNT(*) FILTER (WHERE status = 'verified') as verified_count,
    COUNT(*) FILTER (WHERE status = 'suspicious') as suspicious_count,
    COUNT(*) FILTER (WHERE status = 'scam') as scam_count,
    AVG(risk_score) as avg_risk_score
FROM domains 
WHERE country IS NOT NULL
GROUP BY country 
ORDER BY total_count DESC;

-- name: GetDomainWithModerationStatus :one
SELECT 
    d.*,
    pm.status as moderation_status,
    pm.priority as moderation_priority,
    pm.assigned_to,
    pm.submitted_at as moderation_submitted_at,
    pm.resolved_at as moderation_resolved_at
FROM domains d
LEFT JOIN pending_moderation pm ON d.domain = pm.domain
WHERE d.domain = $1;

-- name: GetSuspiciousDomainsWithModeration :many
SELECT 
    d.*,
    pm.status as moderation_status,
    pm.priority,
    pm.assigned_to,
    pm.submitted_at as moderation_submitted_at
FROM domains d
LEFT JOIN pending_moderation pm ON d.domain = pm.domain
WHERE d.status = 'suspicious'
ORDER BY d.risk_score DESC, pm.priority ASC
LIMIT $1 OFFSET $2;

-- name: GetDomainsByRiskRangeWithModeration :many
SELECT 
    d.*,
    pm.status as moderation_status,
    pm.assigned_to
FROM domains d
LEFT JOIN pending_moderation pm ON d.domain = pm.domain
WHERE d.risk_score BETWEEN $1 AND $2
ORDER BY d.risk_score DESC
LIMIT $3 OFFSET $4;

-- name: GetComprehensiveDomainStats :one
SELECT 
    -- Domain stats
    COUNT(d.*) as total_domains,
    COUNT(*) FILTER (WHERE d.status = 'verified') as verified_domains,
    COUNT(*) FILTER (WHERE d.status = 'suspicious') as suspicious_domains,
    COUNT(*) FILTER (WHERE d.status = 'scam') as scam_domains,
    AVG(d.risk_score) as avg_risk_score,
    
    -- Moderation stats
    COUNT(pm.*) as total_moderation_tasks,
    COUNT(*) FILTER (WHERE pm.status = 'pending') as pending_moderation,
    COUNT(*) FILTER (WHERE pm.status = 'in_progress') as in_progress_moderation,
    
    -- Performance stats
    COUNT(*) FILTER (WHERE d.expires_at < CURRENT_TIMESTAMP AND d.status = 'verified') as expired_domains,
    COUNT(*) FILTER (WHERE d.last_check_at < CURRENT_TIMESTAMP - INTERVAL '7 days') as domains_need_recheck
FROM domains d
LEFT JOIN pending_moderation pm ON d.domain = pm.domain;

-- name: SearchDomainsWithModeration :many
SELECT 
    d.*,
    pm.status as moderation_status,
    pm.priority,
    pm.assigned_to
FROM domains d
LEFT JOIN pending_moderation pm ON d.domain = pm.domain
WHERE d.domain ILIKE $1 
   OR d.company_name ILIKE $1
ORDER BY d.risk_score DESC
LIMIT $2 OFFSET $3;

-- name: GetDashboardData :one
SELECT 
    json_build_object(
        'total_domains', COUNT(d.*),
        'verified_domains', COUNT(*) FILTER (WHERE d.status = 'verified'),
        'suspicious_domains', COUNT(*) FILTER (WHERE d.status = 'suspicious'),
        'scam_domains', COUNT(*) FILTER (WHERE d.status = 'scam'),
        'avg_risk_score', ROUND(AVG(d.risk_score), 2),
        'pending_moderation', COUNT(*) FILTER (WHERE pm.status = 'pending'),
        'in_progress_moderation', COUNT(*) FILTER (WHERE pm.status = 'in_progress'),
        'high_risk_domains', COUNT(*) FILTER (WHERE d.risk_score > 70 AND d.status != 'scam'),
        'expired_domains', COUNT(*) FILTER (WHERE d.expires_at < CURRENT_TIMESTAMP AND d.status = 'verified'),
        'recent_scams', COUNT(*) FILTER (WHERE d.status = 'scam' AND d.updated_at > CURRENT_TIMESTAMP - INTERVAL '24 hours')
    ) as dashboard_stats
FROM domains d
LEFT JOIN pending_moderation pm ON d.domain = pm.domain;

-- name: GetDomainStats :one
SELECT 
    COUNT(*) as total_domains,
    COUNT(*) FILTER (WHERE status = 'verified') as verified_count,
    COUNT(*) FILTER (WHERE status = 'suspicious') as suspicious_count,
    COUNT(*) FILTER (WHERE status = 'scam') as scam_count,
    AVG(risk_score) as avg_risk_score,
    COUNT(*) FILTER (WHERE expires_at < CURRENT_TIMESTAMP AND status = 'verified') as expired_count
FROM domains;

-- name: GetModeratorStats :many
SELECT 
    assigned_to,
    COUNT(*) as total_assigned,
    COUNT(*) FILTER (WHERE status = 'approved') as approved_count,
    COUNT(*) FILTER (WHERE status = 'rejected') as rejected_count,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_count,
    AVG(EXTRACT(EPOCH FROM (resolved_at - submitted_at))/3600) as avg_resolution_time_hours
FROM pending_moderation 
WHERE assigned_to IS NOT NULL
GROUP BY assigned_to 
ORDER BY total_assigned DESC;

-- name: GetRecentlyUpdatedDomains :many
SELECT * FROM domains 
WHERE updated_at > $1 
ORDER BY updated_at DESC 
LIMIT $2;

