package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const trackingQueryTimeout = 5 * time.Second

type TrackingStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewTrackingStore(db *sql.DB, log *zap.SugaredLogger) *TrackingStore {
	return &TrackingStore{
		db:  db,
		log: log,
	}
}

func (s *TrackingStore) GetTripCurrentLocation(
	ctx context.Context,
	tripID string,
) (*TripCurrentLocation, error) {

	const query = `
		SELECT
			trip_id,
			driver_id,
			vehicle_id,
			latitude,
			longitude,
			recorded_at,
			connection
		FROM trip_tracking_current
		WHERE trip_id = $1
	`

	s.log.Debugw("fetching current trip location from database", "trip_id", tripID)

	ctx, cancel := context.WithTimeout(ctx, trackingQueryTimeout)
	defer cancel()

	var location TripCurrentLocation

	err := s.db.QueryRowContext(ctx, query, tripID).Scan(
		&location.TripID,
		&location.DriverID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.RecordedAt,
		&location.Connection,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warnw("current trip location not found", "trip_id", tripID)
			return nil, ErrNotFound
		}

		s.log.Errorw("failed to fetch current trip location", "trip_id", tripID, "err", err)
		return nil, err
	}

	s.log.Debugw("current trip location fetched successfully", "trip_id", tripID)

	return &location, nil
}

func (s *TrackingStore) GetTripLocationHistory(
	ctx context.Context,
	tripID string,
	limit int32,
) ([]*TripLocationHistoryPoint, error) {

	const defaultLimit int32 = 100

	const query = `
		SELECT
			latitude,
			longitude,
			recorded_at
		FROM (
			SELECT
				latitude,
				longitude,
				recorded_at
			FROM trip_tracking_history
			WHERE trip_id = $1
			ORDER BY recorded_at DESC
			LIMIT $2
		) AS recent_points
		ORDER BY recorded_at ASC
	`

	if limit <= 0 {
		limit = defaultLimit
	}

	s.log.Debugw(
		"fetching trip location history from database",
		"trip_id", tripID,
		"limit", limit,
	)

	ctx, cancel := context.WithTimeout(ctx, trackingQueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, tripID, limit)
	if err != nil {
		s.log.Errorw("failed to query trip location history", "trip_id", tripID, "err", err)
		return nil, err
	}
	defer rows.Close()

	points := make([]*TripLocationHistoryPoint, 0)

	for rows.Next() {
		var point TripLocationHistoryPoint

		if err := rows.Scan(
			&point.Latitude,
			&point.Longitude,
			&point.RecordedAt,
		); err != nil {
			s.log.Errorw("failed to scan trip location history point", "trip_id", tripID, "err", err)
			return nil, err
		}

		points = append(points, &point)
	}

	if err := rows.Err(); err != nil {
		s.log.Errorw("trip location history row iteration failed", "trip_id", tripID, "err", err)
		return nil, err
	}

	s.log.Debugw(
		"trip location history fetched successfully",
		"trip_id", tripID,
		"count", len(points),
	)

	return points, nil
}

func (s *TrackingStore) GetTripGeometry(
	ctx context.Context,
	tripID string,
) (*TripGeometry, error) {

	const query = `
		SELECT
			latitude,
			longitude
		FROM trip_tracking_history
		WHERE trip_id = $1
		ORDER BY recorded_at ASC
	`

	s.log.Debugw("fetching trip geometry from database", "trip_id", tripID)

	ctx, cancel := context.WithTimeout(ctx, trackingQueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, tripID)
	if err != nil {
		s.log.Errorw("failed to query trip geometry", "trip_id", tripID, "err", err)
		return nil, err
	}
	defer rows.Close()

	actualGeometry := make([]Coordinate, 0)

	for rows.Next() {
		var point Coordinate

		if err := rows.Scan(&point.Latitude, &point.Longitude); err != nil {
			s.log.Errorw("failed to scan trip geometry point", "trip_id", tripID, "err", err)
			return nil, err
		}

		actualGeometry = append(actualGeometry, point)
	}

	if err := rows.Err(); err != nil {
		s.log.Errorw("trip geometry row iteration failed", "trip_id", tripID, "err", err)
		return nil, err
	}

	if len(actualGeometry) == 0 {
		current, err := s.GetTripCurrentLocation(ctx, tripID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				s.log.Warnw("no stored geometry or current location found for trip", "trip_id", tripID)
				return &TripGeometry{
					TripID:          tripID,
					PlannedGeometry: []Coordinate{},
					ActualGeometry:  []Coordinate{},
				}, nil
			}

			return nil, err
		}

		actualGeometry = append(actualGeometry, Coordinate{
			Latitude:  current.Latitude,
			Longitude: current.Longitude,
		})
	}

	s.log.Debugw(
		"trip geometry fetched successfully",
		"trip_id", tripID,
		"actual_points", len(actualGeometry),
	)

	return &TripGeometry{
		TripID:          tripID,
		PlannedGeometry: []Coordinate{},
		ActualGeometry:  actualGeometry,
	}, nil
}

func (s *TrackingStore) StoreLocationUpdate(
	ctx context.Context,
	update TrackingLocationUpdate,
) error {

	const insertHistoryQuery = `
		INSERT INTO trip_tracking_history (
			tracking_point_id,
			trip_id,
			driver_id,
			vehicle_id,
			latitude,
			longitude,
			recorded_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	const upsertCurrentQuery = `
		INSERT INTO trip_tracking_current (
			trip_id,
			driver_id,
			vehicle_id,
			latitude,
			longitude,
			connection,
			recorded_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (trip_id) DO UPDATE SET
			driver_id = EXCLUDED.driver_id,
			vehicle_id = EXCLUDED.vehicle_id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			connection = EXCLUDED.connection,
			recorded_at = EXCLUDED.recorded_at,
			updated_at = NOW()
	`

	s.log.Debugw(
		"storing tracking location update",
		"trip_id", update.TripID,
		"driver_id", update.DriverID,
		"vehicle_id", update.VehicleID,
		"sequence", update.Sequence,
	)

	ctx, cancel := context.WithTimeout(ctx, trackingQueryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Errorw("failed to begin tracking transaction", "trip_id", update.TripID, "err", err)
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		insertHistoryQuery,
		uuid.NewString(),
		update.TripID,
		update.DriverID,
		update.VehicleID,
		update.Latitude,
		update.Longitude,
		update.RecordedAt,
	); err != nil {
		s.log.Errorw("failed to insert tracking history point", "trip_id", update.TripID, "err", err)
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		upsertCurrentQuery,
		update.TripID,
		update.DriverID,
		update.VehicleID,
		update.Latitude,
		update.Longitude,
		"streaming",
		update.RecordedAt,
	); err != nil {
		s.log.Errorw("failed to upsert current tracking location", "trip_id", update.TripID, "err", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		s.log.Errorw("failed to commit tracking transaction", "trip_id", update.TripID, "err", err)
		return err
	}

	s.log.Debugw(
		"tracking location update stored successfully",
		"trip_id", update.TripID,
		"driver_id", update.DriverID,
		"vehicle_id", update.VehicleID,
		"sequence", update.Sequence,
	)

	return nil
}
