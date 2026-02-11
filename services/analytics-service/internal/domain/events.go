package domain

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"go.uber.org/zap"
)

type EventStore struct {
	db  driver.Conn
	log *zap.SugaredLogger
}

func NewEventsStore(db driver.Conn, log *zap.SugaredLogger) *EventStore {
	return &EventStore{
		db:  db,
		log: log,
	}
}

func (s *EventStore) CountVehiclesInTransit(ctx context.Context) (uint64, error) {
	var count uint64

	query := `
		SELECT count()
		FROM analytics_events
		WHERE event_type = 'vehicle_in_transit';
	`

	s.log.Infow("counting vehicles in transit")

	if err := s.db.QueryRow(ctx, query).Scan(&count); err != nil {
		s.log.Errorw("failed to count vehicles in transit", "err", err)
		return 0, err
	}

	s.log.Infow("vehicles in transit counted successfully",
		"count", count,
	)

	return count, nil
}

func (s *EventStore) CountTripsToday(ctx context.Context) (uint64, error) {
	var count uint64

	startOfDay := time.Now().Truncate(24 * time.Hour)

	query := `
		SELECT count()
		FROM analytics_events
		WHERE event_type = 'trip_completed'
		  AND event_time >= ?
	`

	s.log.Infow("counting trips today",
		"start_of_day", startOfDay,
	)

	if err := s.db.QueryRow(
		ctx,
		query,
		startOfDay,
	).Scan(&count); err != nil {
		s.log.Errorw("failed to count trips today",
			"start_of_day", startOfDay,
			"err", err,
		)
		return 0, err
	}

	s.log.Infow("trips today counted successfully",
		"count", count,
	)

	return count, nil
}

func (s *EventStore) InsertIdentityUserRegistered(
	ctx context.Context,
	eventTime time.Time,
	userID string,
	email string,
) error {

	query := `
		INSERT INTO analytics_events
		(event_time, service, event_type, route, status_code, latency_ms, user_id, request_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	s.log.Infow("inserting identity user registered event",
		"user_id", userID,
		"email", email,
		"event_time", eventTime,
	)

	err := s.db.Exec(
		ctx,
		query,
		eventTime,
		"identity-service",
		"user_registered",
		"",        // route not applicable
		uint16(0), // status_code not applicable
		uint32(0), // latency_ms not applicable
		userID,
		userID, // using userID as request_id for now (can improve later)
	)

	if err != nil {
		s.log.Errorw("failed to insert identity user registered event",
			"user_id", userID,
			"err", err,
		)
		return err
	}

	s.log.Infow("identity user registered event inserted successfully",
		"user_id", userID,
	)

	return nil
}
