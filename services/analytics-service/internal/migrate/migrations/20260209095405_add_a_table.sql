-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS analytics_events
(
    event_time   DateTime,
    service      LowCardinality(String),
    event_type   LowCardinality(String),
    route        String,
    status_code  UInt16,
    latency_ms   UInt32,
    user_id      Nullable(String),
    request_id   String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(event_time)
ORDER BY (service, event_type, event_time)
SETTINGS index_granularity = 8192;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists post;
-- +goose StatementEnd
