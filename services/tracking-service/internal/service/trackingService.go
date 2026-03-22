package service

import (
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
