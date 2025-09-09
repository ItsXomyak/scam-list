-- Down миграция для отката изменений

-- Удаляем триггеры
DROP TRIGGER IF EXISTS trigger_resolve_moderation_task ON domains;
DROP TRIGGER IF EXISTS trigger_create_moderation_task ON domains;
DROP TRIGGER IF EXISTS trigger_auto_set_verified_at ON domains;
DROP TRIGGER IF EXISTS trigger_domains_auto_status_insert ON domains;
DROP TRIGGER IF EXISTS trigger_domains_auto_status ON domains;
DROP TRIGGER IF EXISTS trigger_domains_updated_at ON domains;

-- Удаляем функции
DROP FUNCTION IF EXISTS resolve_moderation_task();
DROP FUNCTION IF EXISTS create_moderation_task();
DROP FUNCTION IF EXISTS auto_set_verified_at();
DROP FUNCTION IF EXISTS auto_set_status_on_insert();
DROP FUNCTION IF EXISTS auto_update_status_by_risk_score();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем таблицу pending_moderation
DROP TABLE IF EXISTS pending_moderation;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_domains_verified_by;
DROP INDEX IF EXISTS idx_domains_scam_type;
DROP INDEX IF EXISTS idx_domains_metadata_gin;
DROP INDEX IF EXISTS idx_domains_reasons_gin;
DROP INDEX IF EXISTS idx_domains_scam_sources_gin;
DROP INDEX IF EXISTS idx_domains_expires_at;
DROP INDEX IF EXISTS idx_domains_last_check_at;
DROP INDEX IF EXISTS idx_domains_updated_at;
DROP INDEX IF EXISTS idx_domains_created_at;
DROP INDEX IF EXISTS idx_domains_status_risk;
DROP INDEX IF EXISTS idx_domains_country;
DROP INDEX IF EXISTS idx_domains_risk_score;
DROP INDEX IF EXISTS idx_domains_status;

-- Удаляем таблицу domains
DROP TABLE IF EXISTS domains;