package service

import (
	"github.com/ThatSneakyCoder/RoutePulse/services/fleet-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/fleet-service/internal/infrastructure/events"
	"go.uber.org/zap"
)

type FleetService struct {
	log          *zap.SugaredLogger
	repo         domain.FleetRepository
	rmqPublisher *events.FleetEventPublisher
}

func NewFleetService(repo domain.FleetRepository, log *zap.SugaredLogger, rmq *events.FleetEventPublisher) *FleetService {
	return &FleetService{
		repo:         repo,
		log:          log,
		rmqPublisher: rmq,
	}
}
