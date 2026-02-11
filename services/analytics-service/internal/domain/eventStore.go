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

func (s *EventStore) InsertEvent(
	ctx context.Context,
	event *AnalyticsEvent,
) error {

	query := `
		INSERT INTO analytics_events
		(event_time, service, event_type, route, status_code, latency_ms, user_id, request_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	s.log.Infow("inserting analytics event",
		"event_type", event.EventType,
		"user_id", event.UserID,
	)

	err := s.db.Exec(
		ctx,
		query,
		event.EventTime,
		event.Service,
		event.EventType,
		event.Route,
		event.StatusCode,
		event.LatencyMs,
		event.UserID,
		event.RequestID,
	)

	if err != nil {
		s.log.Errorw("failed to insert analytics event",
			"event_type", event.EventType,
			"user_id", event.UserID,
			"err", err,
		)
		return err
	}

	return nil
}
