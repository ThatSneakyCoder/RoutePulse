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
	(
		event_time,
		org_id,
		user_id,
		event_type,
		entity_type,
		entity_id,
		service,
		route,
		status_code,
		latency_ms,
		request_id
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	s.log.Infow("inserting analytics event",
		"event_type", event.EventType,
		"user_id", event.UserID,
		"org_id", event.OrgID,
	)

	err := s.db.Exec(
		ctx,
		query,
		event.EventTime,
		event.OrgID,
		event.UserID,
		event.EventType,
		event.EntityType,
		event.EntityID,
		event.Service,
		event.Route,
		event.StatusCode,
		event.LatencyMs,
		event.RequestID,
	)

	if err != nil {
		s.log.Errorw("failed to insert analytics event",
			"event_type", event.EventType,
			"user_id", event.UserID,
			"org_id", event.OrgID,
			"err", err,
		)
		return err
	}

	s.log.Infow("analytics event inserted successfully",
		"event_type", event.EventType,
		"user_id", event.UserID,
	)

	return nil
}

func (s *EventStore) GetRecentActivity(ctx context.Context) ([]*Event, error) {

	query := `
		SELECT event_type, user_id, org_id, service, event_time
		FROM analytics_events
		ORDER BY event_time DESC
		LIMIT 20
	`

	s.log.Infow("fetching recent activity")

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		s.log.Errorw("failed to query recent activity", "err", err)
		return nil, err
	}
	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var e Event

		if err := rows.Scan(
			&e.EventType,
			&e.UserID,
			&e.OrgID,
			&e.Service,
			&e.EventTime,
		); err != nil {
			s.log.Errorw("failed to scan event", "err", err)
			return nil, err
		}

		events = append(events, &e)
	}

	return events, nil
}

func (s *EventStore) CountActiveUsersToday(ctx context.Context) (uint64, error) {
	var count uint64

	query := `
		SELECT count(DISTINCT user_id)
		FROM analytics_events
		WHERE event_type = 'identity.user.logged_in'
		  AND event_time >= today()
	`

	s.log.Infow("counting active users today")

	if err := s.db.QueryRow(ctx, query).Scan(&count); err != nil {
		s.log.Errorw("failed to count active users today", "err", err)
		return 0, err
	}

	return count, nil
}

func (s *EventStore) CountTotalMembers(ctx context.Context) (uint64, error) {
	var count uint64

	query := `
		SELECT count()
		FROM analytics_events
		WHERE event_type = 'organization.member_added'
	`

	s.log.Infow("counting total members")

	if err := s.db.QueryRow(ctx, query).Scan(&count); err != nil {
		s.log.Errorw("failed to count total members", "err", err)
		return 0, err
	}

	return count, nil
}