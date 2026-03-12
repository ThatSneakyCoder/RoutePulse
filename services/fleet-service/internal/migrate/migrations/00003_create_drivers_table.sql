-- +goose Up
-- +goose StatementBegin

CREATE TABLE drivers (
    driver_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,

    vehicle_id UUID,

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_drivers_organization
ON drivers (organization_id);

CREATE INDEX idx_drivers_vehicle
ON drivers (vehicle_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS drivers;
-- +goose StatementEnd
