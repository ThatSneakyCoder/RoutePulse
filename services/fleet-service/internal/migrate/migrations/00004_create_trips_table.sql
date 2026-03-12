-- +goose Up
-- +goose StatementBegin
CREATE TABLE trips (
    trip_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    vehicle_id UUID NOT NULL,
    driver_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'created',
    start_time TIMESTAMP NULL,
    end_time TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_trips_org ON trips(organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trips;
-- +goose StatementEnd
