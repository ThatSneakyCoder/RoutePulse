-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS analytics_events
(
    event_time   DateTime,

    -- multi-tenant
    org_id       String,
    user_id      Nullable(String),

    -- event classification
    event_type   LowCardinality(String),
    entity_type  LowCardinality(String),
    entity_id    Nullable(String),

    -- request metadata
    service      LowCardinality(String),
    route        String,
    status_code  UInt16,
    latency_ms   UInt32,

    -- tracing
    request_id   String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(event_time)
ORDER BY (org_id, event_type, event_time)
SETTINGS index_granularity = 8192;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS analytics_events;
-- +goose StatementEnd