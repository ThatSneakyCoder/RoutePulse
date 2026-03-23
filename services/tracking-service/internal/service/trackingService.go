package service

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/infrastructure/events"
	"go.uber.org/zap"
)

type TrackingService struct {
	log          *zap.SugaredLogger
	repo         domain.TrackingRepository
	rmqPublisher *events.TrackingEventPublisher
}

func NewTrackingService(repo domain.TrackingRepository, log *zap.SugaredLogger, rmq *events.TrackingEventPublisher) *TrackingService {
	return &TrackingService{
		repo:         repo,
		log:          log,
		rmqPublisher: rmq,
	}
}

func (s *TrackingService) GetTripCurrentLocation(
	ctx context.Context,
	tripID string,
) (*domain.TripCurrentLocation, error) {

	s.log.Infow(
		"get trip current location request received",
		"trip_id", tripID,
	)

	if s.repo == nil {
		s.log.Warnw(
			"tracking repository not configured; returning empty current location",
			"trip_id", tripID,
		)

		return &domain.TripCurrentLocation{
			TripID:     tripID,
			Connection: "pending-storage",
		}, nil
	}

	location, err := s.repo.GetTripCurrentLocation(ctx, tripID)
	if err != nil {
		s.log.Errorw(
			"failed to get trip current location",
			"trip_id", tripID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"trip current location fetched successfully",
		"trip_id", tripID,
	)

	return location, nil
}

func (s *TrackingService) GetTripLocationHistory(
	ctx context.Context,
	tripID string,
	limit int32,
) ([]*domain.TripLocationHistoryPoint, error) {

	s.log.Infow(
		"get trip location history request received",
		"trip_id", tripID,
		"limit", limit,
	)

	if s.repo == nil {
		s.log.Warnw(
			"tracking repository not configured; returning empty trip history",
			"trip_id", tripID,
		)

		return []*domain.TripLocationHistoryPoint{}, nil
	}

	points, err := s.repo.GetTripLocationHistory(ctx, tripID, limit)
	if err != nil {
		s.log.Errorw(
			"failed to get trip location history",
			"trip_id", tripID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"trip location history fetched successfully",
		"trip_id", tripID,
		"count", len(points),
	)

	return points, nil
}

func (s *TrackingService) GetTripGeometry(
	ctx context.Context,
	tripID string,
) (*domain.TripGeometry, error) {

	s.log.Infow(
		"get trip geometry request received",
		"trip_id", tripID,
	)

	if s.repo == nil {
		s.log.Warnw(
			"tracking repository not configured; returning empty trip geometry",
			"trip_id", tripID,
		)

		return &domain.TripGeometry{
			TripID:          tripID,
			PlannedGeometry: []domain.Coordinate{},
			ActualGeometry:  []domain.Coordinate{},
		}, nil
	}

	geometry, err := s.repo.GetTripGeometry(ctx, tripID)
	if err != nil {
		s.log.Errorw(
			"failed to get trip geometry",
			"trip_id", tripID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"trip geometry fetched successfully",
		"trip_id", tripID,
	)

	return geometry, nil
}
