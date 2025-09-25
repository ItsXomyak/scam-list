-- Общие таблицы и типы
-- Версия: 1.0

-- Таблица доменов
CREATE TABLE domains (
    domain VARCHAR(253) PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'scam' 
        CHECK (status IN ('verified', 'scam', 'suspicious')),
    
    -- Общая информация
    company_name VARCHAR(255) DEFAULT 'Unknown',
    country CHAR(2) DEFAULT NULL,
    
    -- Данные для scam
    scam_sources VARCHAR(100)[], -- Вместо просто "scam = true" мы знаем какой модуль / источник обнаружил
    scam_type VARCHAR(100) DEFAULT 'other', -- Тип скама, если известен (phishing, malware, fraud, tech_support, lottery, etc.)
    
    -- Данные для verified
    verified_by VARCHAR(100) DEFAULT 'Officers', -- Кто верифицировал, модуль или человек
    verification_method VARCHAR(100) DEFAULT 'manual',
    
    -- Общие поля
    -- risk_score - оценка доверия домена от 0 до 100, где 100 - полностью скамный, 0 - полностью доверенный
    risk_score DECIMAL(5,2) CHECK (risk_score >= 0 AND risk_score <= 100),
    reasons TEXT[],
    metadata JSONB, -- Дополнительные данные в формате JSONB, мол содержать результаты модуей
    
    -- Аудит и timing
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации

-- -- 1. Основные индексы для частых запросов
-- CREATE INDEX idx_domains_status ON domains(status);
-- CREATE INDEX idx_domains_risk_score ON domains(risk_score);
-- CREATE INDEX idx_domains_country ON domains(country) WHERE country IS NOT NULL;

-- -- 2. Составные индексы для комплексных запросов
-- CREATE INDEX idx_domains_status_risk ON domains(status, risk_score);

-- -- 3. Индексы для временных полей (для аналитики и очистки)
-- CREATE INDEX idx_domains_created_at ON domains(created_at);
-- CREATE INDEX idx_domains_updated_at ON domains(updated_at);

-- -- 5. Индексы для массивов (GIN для эффективного поиска)
-- CREATE INDEX idx_domains_scam_sources_gin ON domains USING GIN(scam_sources) 
--     WHERE scam_sources IS NOT NULL;
-- CREATE INDEX idx_domains_reasons_gin ON domains USING GIN(reasons) 
--     WHERE reasons IS NOT NULL;

-- -- 6. Индекс для JSONB метаданных
-- CREATE INDEX idx_domains_metadata_gin ON domains USING GIN(metadata) 
--     WHERE metadata IS NOT NULL;

-- -- 7. Частичные индексы для оптимизации специфических запросов
-- CREATE INDEX idx_domains_scam_type ON domains(scam_type) 
--     WHERE status = 'scam' AND scam_type IS NOT NULL;

-- -- Триггеры и функции

-- CREATE OR REPLACE FUNCTION update_updated_at_column()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = CURRENT_TIMESTAMP;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_domains_updated_at
--     BEFORE UPDATE ON domains
--     FOR EACH ROW
--     EXECUTE FUNCTION update_updated_at_column();


-- CREATE OR REPLACE FUNCTION auto_update_status_by_risk_score()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     -- Проверяем, изменился ли risk_score
--     IF OLD.risk_score IS DISTINCT FROM NEW.risk_score AND NEW.risk_score IS NOT NULL THEN
        
--         -- Определяем статус на основе risk_score
--         CASE
--             WHEN NEW.risk_score >= 0.00 AND NEW.risk_score <= 20.00 THEN
--                 NEW.status = 'verified';
--             WHEN NEW.risk_score >= 20.01 AND NEW.risk_score <= 30.00 THEN
--                 NEW.status = 'verified';
--             WHEN NEW.risk_score >= 30.01 AND NEW.risk_score <= 50.00 THEN
--                 NEW.status = 'suspicious';
--             WHEN NEW.risk_score >= 50.01 AND NEW.risk_score <= 70.00 THEN
--                 NEW.status = 'suspicious';
--             WHEN NEW.risk_score >= 70.01 AND NEW.risk_score <= 90.00 THEN
--                 NEW.status = 'scam';
--             WHEN NEW.risk_score >= 90.01 AND NEW.risk_score <= 100.00 THEN
--                 NEW.status = 'scam';
--         END CASE;
        
--         -- Логирование изменения статуса
--         RAISE NOTICE 'Domain % status auto-updated from % to % based on risk_score %', 
--             NEW.domain, OLD.status, NEW.status, NEW.risk_score;
            
--     END IF;
    
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_domains_auto_status
--     BEFORE UPDATE ON domains
--     FOR EACH ROW
--     EXECUTE FUNCTION auto_update_status_by_risk_score();


-- CREATE OR REPLACE FUNCTION auto_set_status_on_insert()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     -- Если при INSERT указан risk_score, но не указан статус явно
--     IF NEW.risk_score IS NOT NULL AND NEW.status = 'scam' THEN -- проверяем дефолтное значение
--         CASE
--             WHEN NEW.risk_score >= 0.00 AND NEW.risk_score <= 20.00 THEN
--                 NEW.status = 'verified';
--             WHEN NEW.risk_score >= 20.01 AND NEW.risk_score <= 30.00 THEN
--                 NEW.status = 'verified';
--             WHEN NEW.risk_score >= 30.01 AND NEW.risk_score <= 50.00 THEN
--                 NEW.status = 'suspicious';
--             WHEN NEW.risk_score >= 50.01 AND NEW.risk_score <= 70.00 THEN
--                 NEW.status = 'suspicious';
--             WHEN NEW.risk_score >= 70.01 AND NEW.risk_score <= 90.00 THEN
--                 NEW.status = 'scam';
--             WHEN NEW.risk_score >= 90.01 AND NEW.risk_score <= 100.00 THEN
--                 NEW.status = 'scam';
--         END CASE;
--     END IF;
    
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_domains_auto_status_insert
--     BEFORE INSERT ON domains
--     FOR EACH ROW
--     EXECUTE FUNCTION auto_set_status_on_insert();

------------------------------------------------------------------------------------------------------------

-- Таблица доменов на модерации
-- CREATE TABLE pending_moderation (
--     domain VARCHAR(253) PRIMARY KEY,
--     check_id UUID NOT NULL,
--     reasons TEXT[] NOT NULL,
--     source_modules VARCHAR(100)[] NOT NULL,
--     priority INTEGER DEFAULT 5 CHECK (priority >= 1 AND priority <= 10),
--     status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'approved', 'rejected')),
--     assigned_to VARCHAR(100),
--     submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     resolved_at TIMESTAMP WITH TIME ZONE,
--     moderator_notes TEXT,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );