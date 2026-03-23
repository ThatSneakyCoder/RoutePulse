-- +goose Up
-- +goose StatementBegin
ALTER TABLE trips
ADD COLUMN start_location POINT,
ADD COLUMN end_location POINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE trips
DROP COLUMN IF EXISTS start_location,
DROP COLUMN IF EXISTS end_location;
-- +goose StatementEnd
