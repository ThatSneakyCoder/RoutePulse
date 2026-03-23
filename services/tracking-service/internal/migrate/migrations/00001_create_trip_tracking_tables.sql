-- +goose Up
-- +goose StatementBegin
CREATE TABLE trip_tracking_current (
    trip_id UUID PRIMARY KEY,
    driver_id UUID NOT NULL,
    vehicle_id UUID NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    connection TEXT NOT NULL DEFAULT 'offline',
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE trip_tracking_history (
    tracking_point_id UUID PRIMARY KEY,
    trip_id UUID NOT NULL,
    driver_id UUID NOT NULL,
    vehicle_id UUID NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trip_tracking_history_trip_recorded
ON trip_tracking_history (trip_id, recorded_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trip_tracking_history;
DROP TABLE IF EXISTS trip_tracking_current;
-- +goose StatementEnd
