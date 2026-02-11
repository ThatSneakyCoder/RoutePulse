package service

import (
	"context"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/domain"
)

func (s *AnalyticsService) InsertDomainEvent(
	ctx context.Context,
	serviceName string,
	eventType string,
	userID string,
	requestID string,
) error {

	s.log.Infow("handling domain analytics event",
		"service", serviceName,
		"event_type", eventType,
		"user_id", userID,
	)

	analyticsEvent := &domain.AnalyticsEvent{
		EventTime:  time.Now(),
		Service:    serviceName,
		EventType:  eventType,
		UserID:     userID,
		RequestID:  requestID,
		Route:      "",
		StatusCode: 0,
		LatencyMs:  0,
	}

	if err := s.repo.InsertEvent(ctx, analyticsEvent); err != nil {
		s.log.Errorw("failed to insert analytics event",
			"event_type", eventType,
			"user_id", userID,
			"err", err,
		)
		return err
	}

	return nil
}
