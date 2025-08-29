-- name: CreateModerationTask :one
INSERT INTO pending_moderation (
    domain, check_id, reasons, source_modules, priority, status
) VALUES (
    $1, $2, $3, $4, $5, 'pending'
) RETURNING *;

-- name: GetModerationTask :one
SELECT * FROM pending_moderation WHERE domain = $1;

-- name: GetModerationTaskByCheckId :one
SELECT * FROM pending_moderation WHERE check_id = $1;

-- name: GetPendingModerationTasks :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.status = 'pending'
ORDER BY pm.priority ASC, pm.submitted_at ASC
LIMIT $1 OFFSET $2;

-- name: GetModerationTasksByStatus :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.status = $1
ORDER BY pm.submitted_at DESC
LIMIT $2 OFFSET $3;

-- name: GetModerationTasksByModerator :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.assigned_to = $1
ORDER BY pm.submitted_at DESC
LIMIT $2 OFFSET $3;

-- name: AssignModerationTask :one
UPDATE pending_moderation 
SET status = 'in_progress', assigned_to = $2, updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 AND status = 'pending'
RETURNING *;

-- name: UpdateModerationTaskStatus :one
UPDATE pending_moderation 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING *;

-- name: ResolveModerationTask :one
UPDATE pending_moderation 
SET status = $2, resolved_at = CURRENT_TIMESTAMP, moderator_notes = $3, updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING *;

-- name: DeleteModerationTask :exec
DELETE FROM pending_moderation WHERE domain = $1;

-- name: GetActiveModerationTasks :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.status IN ('pending', 'in_progress')
ORDER BY pm.priority ASC, pm.submitted_at ASC;

-- name: GetModerationTasksByPriority :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.priority <= $1 AND pm.status = 'pending'
ORDER BY pm.priority ASC, pm.submitted_at ASC
LIMIT $2 OFFSET $3;

-- name: GetOverdueModerationTasks :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.status = 'in_progress' 
AND pm.submitted_at < $1
ORDER BY pm.submitted_at ASC;

-- name: UpdateModerationTaskPriority :one
UPDATE pending_moderation 
SET priority = $2, updated_at = CURRENT_TIMESTAMP
WHERE domain = $1 
RETURNING *;

-- name: GetModerationStats :one
SELECT 
    COUNT(*) as total_tasks,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_count,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_count,
    COUNT(*) FILTER (WHERE status = 'approved') as approved_count,
    COUNT(*) FILTER (WHERE status = 'rejected') as rejected_count,
    AVG(EXTRACT(EPOCH FROM (COALESCE(resolved_at, CURRENT_TIMESTAMP) - submitted_at))/3600) as avg_resolution_time_hours
FROM pending_moderation;



-- name: GetModerationTasksBySourceModule :many
SELECT 
    UNNEST(source_modules) as source_module,
    COUNT(*) as task_count,
    AVG(priority) as avg_priority
FROM pending_moderation 
GROUP BY source_module 
ORDER BY task_count DESC;

-- name: GetRecentModerationActivity :many
SELECT pm.*, d.risk_score, d.company_name, d.country
FROM pending_moderation pm
JOIN domains d ON pm.domain = d.domain
WHERE pm.resolved_at > $1 OR pm.submitted_at > $1
ORDER BY COALESCE(pm.resolved_at, pm.submitted_at) DESC
LIMIT $2;

-- name: GetHighRiskDomains :many
SELECT * FROM domains 
WHERE risk_score >= $1 
AND status = 'suspicious' 
ORDER BY risk_score DESC 
LIMIT $2;

-- name: SearchDomainsByPattern :many
SELECT * FROM domains 
WHERE domain ILIKE $1 
ORDER BY risk_score DESC 
LIMIT $2 OFFSET $3;