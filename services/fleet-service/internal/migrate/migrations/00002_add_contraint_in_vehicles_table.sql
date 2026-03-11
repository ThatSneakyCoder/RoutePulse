-- +goose Up
-- +goose StatementBegin
ALTER TABLE vehicles
ADD CONSTRAINT unique_org_plate
UNIQUE (organization_id, plate_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vehicles
DROP CONSTRAINT IF EXISTS unique_org_plate;
-- +goose StatementEnd
