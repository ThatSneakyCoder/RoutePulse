-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS analytics_events;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- no rollback (data loss)
-- +goose StatementEnd