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
      verified_at TIMESTAMP WITH TIME ZONE,
      verified_by VARCHAR(100) DEFAULT 'Officers', -- Кто верифицировал, модуль или человек
      verification_method VARCHAR(100) DEFAULT 'manual',
      expires_at TIMESTAMP WITH TIME ZONE,
      
      -- Общие поля
      -- risk_score - оценка доверия домена от 0 до 100, где 100 - полностью скамный, 0 - полностью доверенный
      risk_score DECIMAL(5,2) CHECK (risk_score >= 0 AND risk_score <= 100),
      reasons TEXT[],
      metadata JSONB, -- Дополнительные данные в формате JSONB, мол содержать результаты модуей
      
      -- Аудит и timing
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      last_check_at TIMESTAMP WITH TIME ZONE  
  );

-- Индексы для оптимизации

  -- 1. Основные индексы для частых запросов
  CREATE INDEX idx_domains_status ON domains(status);
  CREATE INDEX idx_domains_risk_score ON domains(risk_score);
  CREATE INDEX idx_domains_country ON domains(country) WHERE country IS NOT NULL;

  -- 2. Составные индексы для комплексных запросов
  CREATE INDEX idx_domains_status_risk ON domains(status, risk_score);

  -- 3. Индексы для временных полей (для аналитики и очистки)
  CREATE INDEX idx_domains_created_at ON domains(created_at);
  CREATE INDEX idx_domains_updated_at ON domains(updated_at);
  CREATE INDEX idx_domains_last_check_at ON domains(last_check_at) WHERE last_check_at IS NOT NULL;

  -- 4. Индекс для поиска по истечению верификации
  CREATE INDEX idx_domains_expires_at ON domains(expires_at) 
      WHERE expires_at IS NOT NULL AND status = 'verified';

  -- 5. Индексы для массивов (GIN для эффективного поиска)
  CREATE INDEX idx_domains_scam_sources_gin ON domains USING GIN(scam_sources) 
      WHERE scam_sources IS NOT NULL;
  CREATE INDEX idx_domains_reasons_gin ON domains USING GIN(reasons) 
      WHERE reasons IS NOT NULL;

  -- 6. Индекс для JSONB метаданных
  CREATE INDEX idx_domains_metadata_gin ON domains USING GIN(metadata) 
      WHERE metadata IS NOT NULL;

  -- 7. Частичные индексы для оптимизации специфических запросов
  CREATE INDEX idx_domains_scam_type ON domains(scam_type) 
      WHERE status = 'scam' AND scam_type IS NOT NULL;
  CREATE INDEX idx_domains_verified_by ON domains(verified_by) 
      WHERE status = 'verified' AND verified_by IS NOT NULL;

-- Триггеры и функции

    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER trigger_domains_updated_at
        BEFORE UPDATE ON domains
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();


    CREATE OR REPLACE FUNCTION auto_update_status_by_risk_score()
    RETURNS TRIGGER AS $$
    BEGIN
        -- Проверяем, изменился ли risk_score
        IF OLD.risk_score IS DISTINCT FROM NEW.risk_score AND NEW.risk_score IS NOT NULL THEN
            
            -- Определяем статус на основе risk_score
            CASE
                WHEN NEW.risk_score >= 0.00 AND NEW.risk_score <= 20.00 THEN
                    NEW.status = 'verified';
                WHEN NEW.risk_score >= 20.01 AND NEW.risk_score <= 30.00 THEN
                    NEW.status = 'verified';
                WHEN NEW.risk_score >= 30.01 AND NEW.risk_score <= 50.00 THEN
                    NEW.status = 'suspicious';
                WHEN NEW.risk_score >= 50.01 AND NEW.risk_score <= 70.00 THEN
                    NEW.status = 'suspicious';
                WHEN NEW.risk_score >= 70.01 AND NEW.risk_score <= 90.00 THEN
                    NEW.status = 'scam';
                WHEN NEW.risk_score >= 90.01 AND NEW.risk_score <= 100.00 THEN
                    NEW.status = 'scam';
            END CASE;
            
            -- Логирование изменения статуса
            RAISE NOTICE 'Domain % status auto-updated from % to % based on risk_score %', 
                NEW.domain, OLD.status, NEW.status, NEW.risk_score;
                
        END IF;
        
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER trigger_domains_auto_status
        BEFORE UPDATE ON domains
        FOR EACH ROW
        EXECUTE FUNCTION auto_update_status_by_risk_score();


    CREATE OR REPLACE FUNCTION auto_set_status_on_insert()
    RETURNS TRIGGER AS $$
    BEGIN
        -- Если при INSERT указан risk_score, но не указан статус явно
        IF NEW.risk_score IS NOT NULL AND NEW.status = 'scam' THEN -- проверяем дефолтное значение
            CASE
                WHEN NEW.risk_score >= 0.00 AND NEW.risk_score <= 20.00 THEN
                    NEW.status = 'verified';
                WHEN NEW.risk_score >= 20.01 AND NEW.risk_score <= 30.00 THEN
                    NEW.status = 'verified';
                WHEN NEW.risk_score >= 30.01 AND NEW.risk_score <= 50.00 THEN
                    NEW.status = 'suspicious';
                WHEN NEW.risk_score >= 50.01 AND NEW.risk_score <= 70.00 THEN
                    NEW.status = 'suspicious';
                WHEN NEW.risk_score >= 70.01 AND NEW.risk_score <= 90.00 THEN
                    NEW.status = 'scam';
                WHEN NEW.risk_score >= 90.01 AND NEW.risk_score <= 100.00 THEN
                    NEW.status = 'scam';
            END CASE;
        END IF;
        
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER trigger_domains_auto_status_insert
        BEFORE INSERT ON domains
        FOR EACH ROW
        EXECUTE FUNCTION auto_set_status_on_insert();

    
    -- Автоматическая установка verified_at при смене статуса на verified
    CREATE OR REPLACE FUNCTION auto_set_verified_at()
      RETURNS TRIGGER AS $$
      BEGIN
          -- Проверяем, что статус верифицированный И verified_at не задан явно
          IF NEW.status = 'verified' AND NEW.verified_at IS NULL THEN
              NEW.verified_at = CURRENT_TIMESTAMP;
          END IF;
          
          RETURN NEW;
      END;
      $$ LANGUAGE plpgsql;

    CREATE TRIGGER trigger_auto_set_verified_at
        BEFORE INSERT ON domains
        FOR EACH ROW
        EXECUTE FUNCTION auto_set_verified_at();

------------------------------------------------------------------------------------------------------------

-- Триггеры для создания и завершения задач модерации
  -- Триггер для создания задачи на модерацию при изменении статуса на suspicious
  CREATE OR REPLACE FUNCTION create_moderation_task()
  RETURNS TRIGGER AS $$
  BEGIN
      -- Проверяем, изменился ли статус на 'suspicious'
      IF (OLD IS NULL OR OLD.status != 'suspicious') AND NEW.status = 'suspicious' THEN
          
          -- Проверяем, нет ли уже активной задачи для этого домена
          IF NOT EXISTS (
              SELECT 1 FROM pending_moderation 
              WHERE domain = NEW.domain 
              AND status IN ('pending', 'in_progress')
          ) THEN
              -- Создаем новую задачу для модератора
              INSERT INTO pending_moderation (
                  domain,
                  check_id,
                  reasons,
                  source_modules,
                  priority,
                  status,
                  submitted_at
              ) VALUES (
                  NEW.domain,
                  gen_random_uuid(), -- Генерируем уникальный ID проверки
                  COALESCE(NEW.reasons, ARRAY['Автоматическая проверка выявила подозрительную активность']),
                  COALESCE(NEW.scam_sources, ARRAY['auto_risk_assessment']),
                  CASE 
                      WHEN NEW.risk_score >= 60 THEN 3  -- Высокий приоритет
                      WHEN NEW.risk_score >= 40 THEN 5  -- Средний приоритет
                      ELSE 7                            -- Низкий приоритет
                  END,
                  'pending',
                  CURRENT_TIMESTAMP
              );
              
              RAISE NOTICE 'Created moderation task for domain: % with priority: %', 
                  NEW.domain, 
                  CASE 
                      WHEN NEW.risk_score >= 60 THEN 3
                      WHEN NEW.risk_score >= 40 THEN 5
                      ELSE 7
                  END;
          END IF;
      END IF;
      
      RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE TRIGGER trigger_create_moderation_task
      AFTER INSERT OR UPDATE ON domains
      FOR EACH ROW
      EXECUTE FUNCTION create_moderation_task();
  ------------------------------------------------------------------------------------------------------------

 -- Срабатывает когда статус домена меняется с 'suspicious' на 'verified' или 'scam'

  CREATE OR REPLACE FUNCTION resolve_moderation_task()
  RETURNS TRIGGER AS $$
  DECLARE
      moderation_record RECORD;
  BEGIN
      -- Проверяем, изменился ли статус с 'suspicious' на финальный
      IF OLD.status = 'suspicious' AND NEW.status IN ('verified', 'scam') THEN
          
          -- Находим активную задачу модерации
          SELECT * INTO moderation_record
          FROM pending_moderation 
          WHERE domain = NEW.domain 
          AND status IN ('pending', 'in_progress')
          ORDER BY submitted_at DESC 
          LIMIT 1;
          
          -- Если задача найдена, завершаем её
          IF FOUND THEN
              UPDATE pending_moderation 
              SET 
                  status = 'resolved',
                  resolved_at = CURRENT_TIMESTAMP,
                  moderator_notes = COALESCE(moderator_notes, '') || 
                      format(' [AUTO] Status changed to %s with risk_score: %s', NEW.status, NEW.risk_score)
              WHERE domain = NEW.domain 
              AND check_id = moderation_record.check_id;
              
              RAISE NOTICE 'Resolved moderation task for domain: % with final status: %', 
                  NEW.domain, NEW.status;
          END IF;
      END IF;
      
      RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE TRIGGER trigger_resolve_moderation_task
      AFTER UPDATE ON domains
      FOR EACH ROW
      EXECUTE FUNCTION resolve_moderation_task();

------------------------------------------------------------------------------------------------------------

-- TODO: доделать функции для модеров: взять задачу в работу, обновить приоритет, добавить заметки, завершить задачу

------------------------------------------------------------------------------------------------------------

-- Таблица доменов на модерации
CREATE TABLE pending_moderation (
    domain VARCHAR(253) PRIMARY KEY,
    check_id UUID NOT NULL,
    reasons TEXT[] NOT NULL,
    source_modules VARCHAR(100)[] NOT NULL,
    priority INTEGER DEFAULT 5 CHECK (priority >= 1 AND priority <= 10),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'approved', 'rejected')),
    assigned_to VARCHAR(100),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    moderator_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);