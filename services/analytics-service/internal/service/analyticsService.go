package service

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/domain"
	"go.uber.org/zap"
)

type AnalyticsService struct {
	repo domain.EventRepository
	log  *zap.SugaredLogger
}

func NewAnalyticsService(repo domain.EventRepository, log *zap.SugaredLogger) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
		log:  log,
	}
}

func (s *AnalyticsService) GetVehiclesInTransit(ctx context.Context) (uint64, error) {
	s.log.Infow("counting vehicles in transit")

	count, err := s.repo.CountVehiclesInTransit(ctx)
	if err != nil {
		s.log.Errorw("failed to get vehicles in transit", "err", err)
		return 0, err
	}

	return count, nil
}

func (s *AnalyticsService) GetTripsToday(ctx context.Context) (uint64, error) {
	s.log.Infow("counting trips today")

	count, err := s.repo.CountTripsToday(ctx)
	if err != nil {
		s.log.Errorw("failed to get trips today", "err", err)
		return 0, err
	}

	return count, nil
}

func (s *AnalyticsService) GetTotalMembers(ctx context.Context) (uint64, error) {
	s.log.Infow("counting total members")

	count, err := s.repo.CountTotalMembers(ctx)
	if err != nil {
		s.log.Errorw("failed to count total members", "err", err)
		return 0, err
	}

	return count, nil
}

func (s *AnalyticsService) GetActiveUsersToday(ctx context.Context) (uint64, error) {
	s.log.Infow("counting active users today")

	count, err := s.repo.CountActiveUsersToday(ctx)
	if err != nil {
		s.log.Errorw("failed to count active users today", "err", err)
		return 0, err
	}

	return count, nil
}

func (s *AnalyticsService) GetRecentActivity(ctx context.Context) ([]*domain.Event, error) {
	s.log.Infow("fetching recent activity")

	events, err := s.repo.GetRecentActivity(ctx)
	if err != nil {
		s.log.Errorw("failed to fetch recent activity", "err", err)
		return nil, err
	}

	return events, nil
}
