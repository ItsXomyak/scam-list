# Общее руководство по таблицам 

**Версия API:** 1.0
**Последнее обновление:** 2025-08-28

## Оглавление

1.  [Обзор системы](#обзор-системы)
2.  [Схема базы данных](#схема-базы-данных)
    *   [Таблица `domains`](#таблица-domains)
    *   [Таблица `pending_moderation`](#таблица-pending_moderation)
    *   [Таблица `domain_checks`](#таблица-domain_checks)
3.  [Бизнес-логика и триггеры](#бизнес-логика-и-триггеры)
    *   [Автоматический статус по risk_score](#автоматический-статус-по-risk_score)
    *   [Модерация доменов](#модерация-доменов)
4.  [Индексы](#индексы)
5.  [Примеры данных](#примеры-данных)
6.  [Типовые сценарии использования](#типовые-сценарии-использования)

---

### Обзор системы

Система предназначена для анализа, категоризации и модерации доменов на основе оценки доверия (`risk_score`). Каждому домену присваивается статус (`verified`, `suspicious`, `scam`) на основе автоматических проверок и, при необходимости, ручной модерации.

**Ключевые понятия:**
*   **`risk_score`** (0.00 - 100.00): Числовая оценка доверия к домену.
    *   `0.00 - 20.00`: Идеальный (Trusted)
    *   `20.01 - 30.00`: Высокое доверие (Verified)
    *   `30.01 - 50.00`: Среднее доверие (Suspicious)
    *   `50.01 - 70.00`: Низкое доверие (Suspicious)
    *   `70.01 - 90.00`: Высокий риск (Scam)
    *   `90.01 - 100.00`: Подтвержденный мошеннический (Scam)
*   **Статус (`status`)**: Категория домена, определяемая автоматически на основе `risk_score`.
*   **Модерация**: Процесс ручной проверки доменов со статусом `suspicious`.

---

### Схема базы данных

#### Таблица `domains`
Основная таблица для хранения информации о доменах.

| Поле | Тип | Ограничения | Описание |
| :--- | :--- | :--- | :--- |
| `domain` | `VARCHAR(253)` | `PRIMARY KEY` | Доменное имя (FQDN). |
| `status` | `VARCHAR(20)` | `NOT NULL DEFAULT 'scam' CHECK IN ('verified', 'scam', 'suspicious')` | Текущий статус домена. |
| `company_name` | `VARCHAR(255)` | `DEFAULT 'Unknown'` | Название компании-владельца. |
| `country` | `CHAR(2)` | `DEFAULT NULL` | Код страны (ISO 3166-1 alpha-2). |
| `scam_sources` | `VARCHAR(100)[]` | - | Массив источников, которые пометили домен как мошеннический. |
| `scam_type` | `VARCHAR(100)` | `DEFAULT 'other'` | Тип мошенничества (phishing, malware, fraud, etc.). |
| `verified_at` | `TIMESTAMPTZ` | - | Время последней верификации. |
| `verified_by` | `VARCHAR(100)` | `DEFAULT 'Officers'` | Источник верификации (модуль или модератор). |
| `verification_method` | `VARCHAR(100)` | `DEFAULT 'manual'` | Метод верификации. |
| `expires_at` | `TIMESTAMPTZ` | - | Срок действия верификации. |
| `risk_score` | `DECIMAL(5,2)` | `CHECK (0 <= risk_score <= 100)`**Ключевое поле.** Оценка риска домена. |
| `reasons` | `TEXT[]` | - | Массив причин для текущего статуса/оценки. |
| `metadata` | `JSONB` | - | Дополнительные данные результатов проверок модулей. |
| `created_at` | `TIMESTAMPTZ` | `DEFAULT CURRENT_TIMESTAMP` | Время создания записи. |
| `updated_at` | `TIMESTAMPTZ` | `DEFAULT CURRENT_TIMESTAMP` | Время последнего обновления записи. |
| `last_check_at` | `TIMESTAMPTZ` | - | Время последней автоматической проверки. |

#### Таблица `pending_moderation`
Очередь доменов, ожидающих ручной модерации.

| Поле | Тип | Описание |
| :--- | :--- | :--- |
| `domain` | `VARCHAR(253)` | `PRIMARY KEY`, ссылка на `domains.domain`. |
| `check_id` | `UUID` | `NOT NULL`, уникальный идентификатор задачи проверки. |
| `reasons` | `TEXT[]` | `NOT NULL`, причины, по которым домен требует модерации. |
| `source_modules` | `VARCHAR(100)[]` | `NOT NULL`, модули, выявившие подозрительную активность. |
| `priority` | `INTEGER` | Приоритет задачи (`1-10`). Высокий risk_score = высокий приоритет (меньшее число). |
| `status` | `VARCHAR(20)` | Статус задачи: `pending`, `in_progress`, `resolved`. |
| `assigned_to` | `VARCHAR(100)` | Идентификатор модератора, взявшего задачу. |
| `submitted_at` | `TIMESTAMPTZ` | Время создания задачи. |
| `resolved_at` | `TIMESTAMPTZ` | Время завершения модерации. |
| `moderator_notes` | `TEXT` | Заметки модератора по результатам проверки. |
| `created_at` | `TIMESTAMPTZ` | Время создания записи. |

---

### Бизнес-логика и триггеры

#### Автоматический статус по risk_score
Система автоматически управляет статусом домена на основе поля `risk_score`.
*   **Триггер `trigger_domains_auto_status`**: При `UPDATE` пересчитывает статус, если изменился `risk_score`.
*   **Триггер `trigger_domains_auto_status_insert`**: При `INSERT` устанавливает статус based on provided `risk_score`.
*   **Триггер `trigger_domains_updated_at`**: Автоматически обновляет `updated_at` при любом изменении записи.

#### Модерация доменов
Для доменов со статусом `suspicious` автоматически создаются задачи для ручной проверки.
*   **Триггер `trigger_create_moderation_task`**: Создает запись в `pending_moderation` при смене статуса на `suspicious`.
*   **Триггер `trigger_resolve_moderation_task`**: Автоматически закрывает активную задачу модерации при смене статуса домена с `suspicious` на `verified` или `scam`.
*   **Триггер `trigger_auto_set_verified_at`**: Автоматически ставит время для verified_at если модератор отправил запись со статусом verified, а само время в запросе verify = NULL

---

### Индексы
Для обеспечения высокой производительности созданы следующие индексы:
*   **Основные**: `status`, `risk_score`, `country`.
*   **Составные**: `(status, risk_score)`, `(status, is_active)`.
*   **Для временных полей**: `created_at`, `updated_at`, `last_check_at`, `expires_at`.
*   **Для массивов и JSONB (GIN)**: `scam_sources`, `reasons`, `metadata`.
*   **Частичные (Partial)**: Для оптимизации запросов к определенным статусам (`scam_type` для `scam`, `verified_by` для `verified`).

---

### Примеры данных

```sql
-- 1. Полностью верифицированный банк
INSERT INTO domains (domain, status, risk_score, company_name, country, verified_at)
VALUES ('kaspi.kz', 'verified', 5.00, 'Kaspi Bank JSC', 'KZ', NOW());

-- 2. Верифицированный, но с признаками для наблюдения
INSERT INTO domains (domain, status, risk_score, company_name, reasons)
VALUES ('new-shop.kz', 'verified', 35.00, 'New Shop LLC', ARRAY['New domain', 'Low traffic']);

-- 3. Подтвержденный скам
INSERT INTO domains (domain, status, risk_score, scam_sources, scam_type, reasons)
VALUES ('scam-site.kz', 'scam', 100.00, ARRAY['VirusTotal', 'InternalDetector'], 'phishing', ARRAY['Phishing kit detected', 'Reported by 50+ users']);

-- 4. Подозрительный домен (попадает в очередь на модерацию)
INSERT INTO domains (domain, status, risk_score, reasons, scam_sources)
VALUES ('fishy-site.kz', 'suspicious', 68.00, ARRAY['High domain age score'], ARRAY['UrlVoid']);
```

### Типовые сценарии использования

1.  **Добавление нового домена после проверки**:
    *   Вставляется запись в `domain_checks`.
    *   Рассчитывается общий `risk_score`.
    *   Вставляется или обновляется запись в `domains`. Если статус установился как `suspicious` — автоматически создается задача в `pending_moderation`.

2.  **Ручная модерация**:
    *   Модератор запрашивает задачи: `SELECT * FROM pending_moderation WHERE status = 'pending' ORDER BY priority ASC`.
    *   Берет задачу в работу: `UPDATE pending_moderation SET status='in_progress', assigned_to=$mod_id WHERE domain=$1`.
    *   После проверки модератор обновляет статус домена: `UPDATE domains SET status='verified', risk_score=15.00, ... WHERE domain=$1`. Триггер автоматически закроет задачу модерации.

3.  **Получение статистики**:
    ```sql
    -- Количество доменов по статусам
    SELECT status, COUNT(*) FROM domains GROUP BY status;

    -- Средний risk_score по странам
    SELECT country, AVG(risk_score) as avg_risk
    FROM domains WHERE country IS NOT NULL GROUP BY country ORDER BY avg_risk DESC;
    ```
### Общая логика модерации доменов
  -- Домен проходит автоматические проверки через разные модули (VirusTotal, URLVoid, собственные алгоритмы)
  -- Система вычисляет общий risk_score от 0 до 100
  -- Если статус стал suspicious → автоматически создается задача в pending_moderation
  -- Приоритет задачи зависит от риска: 60%+ = высокий, 40-60% = средний, меньше 40% = низкий
  -- В задаче указываются причины подозрений и какие модули их выявили
  -- Модератор видит очередь задач, отсортированную по приоритету
  -- Берет задачу в работу (статус in_progress)
  -- Проверяет домен вручную, изучает содержимое, читает отзывы
  -- Принимает финальное решение: verified (безопасен) или scam (мошенник)
  -- Модератор устанавливает финальный статус и risk_score в основной таблице
  -- Задача автоматически помечается как resolved с временной меткой
  -- Сохраняются заметки модератора для истории
  -- Домен получает окончательный статус в основной таблице
  -- История модерации сохраняется для аналитики и контроля качества
  -- Можно отслеживать производительность модераторов и точность автопроверок