-- +goose Up
-- +goose StatementBegin
CREATE TABLE vehicles (
    vehicle_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,

    plate_number TEXT NOT NULL,
    vehicle_type TEXT,
    capacity INTEGER,

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vehicles_organization
ON vehicles (organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vehicles;
-- +goose StatementEnd
