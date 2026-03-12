package service

import (
	"context"
	"errors"

	"github.com/ThatSneakyCoder/RoutePulse/services/fleet-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/fleet-service/internal/infrastructure/events"
	"github.com/google/uuid"
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

func (s *FleetService) CreateVehicle(
	ctx context.Context,
	orgID string,
	plate string,
	vType string,
	cap int32,
) (*domain.Vehicle, error) {

	s.log.Infow(
		"creating vehicle",
		"organization_id", orgID,
		"plate_number", plate,
	)

	var allowedVehicleTypes = map[string]bool{
		"truck": true,
		"van":   true,
		"car":   true,
		"bike":  true,
	}

	if !allowedVehicleTypes[vType] {
		return nil, errors.New("invalid vehicle type")
	}

	v := &domain.Vehicle{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		PlateNumber:    plate,
		VehicleType:    vType,
		Capacity:       cap,
		Status:         "active",
	}

	err := s.repo.CreateVehicle(ctx, v)
	if err != nil {

		s.log.Errorw(
			"failed to create vehicle",
			"error", err,
			"organization_id", orgID,
			"plate_number", plate,
		)

		return nil, err
	}

	s.log.Infow(
		"vehicle created",
		"vehicle_id", v.ID,
		"organization_id", orgID,
	)

	return v, nil
}

func (s *FleetService) ListVehicles(
	ctx context.Context,
	orgID string,
) ([]*domain.Vehicle, error) {

	s.log.Infow(
		"listing vehicles",
		"organization_id", orgID,
	)

	vehicles, err := s.repo.ListVehicles(ctx, orgID)
	if err != nil {

		s.log.Errorw(
			"failed to list vehicles",
			"error", err,
			"organization_id", orgID,
		)

		return nil, err
	}

	s.log.Infow(
		"vehicles listed successfully",
		"organization_id", orgID,
		"count", len(vehicles),
	)

	return vehicles, nil
}
