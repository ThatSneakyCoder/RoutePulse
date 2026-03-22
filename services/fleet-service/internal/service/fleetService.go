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

func (s *FleetService) GetVehicle(
	ctx context.Context,
	vehicleID string,
) (*domain.Vehicle, error) {

	s.log.Infow(
		"fetch vehicle request",
		"vehicle_id", vehicleID,
	)

	v, err := s.repo.GetVehicle(ctx, vehicleID)
	if err != nil {

		if err == domain.ErrNotFound {
			s.log.Warnw(
				"vehicle not found",
				"vehicle_id", vehicleID,
			)
			return nil, err
		}

		s.log.Errorw(
			"failed to fetch vehicle",
			"vehicle_id", vehicleID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"vehicle fetched successfully",
		"vehicle_id", vehicleID,
	)

	return v, nil
}

func (s *FleetService) CreateDriver(
	ctx context.Context,
	orgID string,
	firstName string,
	lastName string,
	vehicleID *string,
) (*domain.Driver, error) {

	s.log.Infow(
		"create driver request",
		"organization_id", orgID,
	)

	driver := &domain.Driver{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		FirstName:      firstName,
		LastName:       lastName,
		VehicleID:      vehicleID,
		Status:         "active",
	}

	err := s.repo.CreateDriver(ctx, driver)
	if err != nil {

		s.log.Errorw(
			"failed to create driver",
			"organization_id", orgID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"driver created successfully",
		"driver_id", driver.ID,
	)

	return driver, nil
}

func (s *FleetService) ListDrivers(
	ctx context.Context,
	orgID string,
) ([]*domain.Driver, error) {

	s.log.Infow(
		"list drivers request",
		"organization_id", orgID,
	)

	drivers, err := s.repo.ListDrivers(ctx, orgID)
	if err != nil {

		s.log.Errorw(
			"failed to list drivers",
			"organization_id", orgID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"drivers fetched successfully",
		"organization_id", orgID,
		"count", len(drivers),
	)

	return drivers, nil
}

func (s *FleetService) CreateTrip(
	ctx context.Context,
	orgID string,
	vehicleID string,
	driverID string,
) (*domain.Trip, error) {

	s.log.Infow(
		"service create trip request received",
		"organization_id", orgID,
		"vehicle_id", vehicleID,
		"driver_id", driverID,
	)

	if err := validateTripOwnership(ctx, s.repo, orgID, vehicleID, driverID); err != nil {

		s.log.Warnw(
			"trip ownership validation failed",
			"organization_id", orgID,
			"vehicle_id", vehicleID,
			"driver_id", driverID,
			"error", err,
		)

		return nil, err
	}

	if err := ensureVehicleAvailable(ctx, s.repo, vehicleID); err != nil {

		s.log.Warnw(
			"vehicle already has active trip",
			"vehicle_id", vehicleID,
			"error", err,
		)

		return nil, err
	}

	if err := ensureDriverAvailable(ctx, s.repo, driverID); err != nil {

		s.log.Warnw(
			"driver already has active trip",
			"driver_id", driverID,
			"error", err,
		)

		return nil, err
	}

	trip := &domain.Trip{
		ID:             uuid.NewString(),
		OrganizationID: orgID,
		VehicleID:      vehicleID,
		DriverID:       driverID,
		Status:         "created",
	}

	err := s.repo.CreateTrip(ctx, trip)
	if err != nil {

		s.log.Errorw(
			"failed to create trip",
			"trip_id", trip.ID,
			"organization_id", orgID,
			"vehicle_id", vehicleID,
			"driver_id", driverID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"trip created successfully",
		"trip_id", trip.ID,
		"organization_id", orgID,
		"vehicle_id", vehicleID,
		"driver_id", driverID,
	)

	return trip, nil
}

func (s *FleetService) ListTrips(
	ctx context.Context,
	orgID string,
) ([]*domain.Trip, error) {

	s.log.Infow(
		"service list trips",
		"organization_id", orgID,
	)

	trips, err := s.repo.ListTrips(ctx, orgID)
	if err != nil {

		s.log.Errorw(
			"service list trips failed",
			"error", err,
		)

		return nil, err
	}

	return trips, nil
}

func (s *FleetService) UpdateVehicle(
	ctx context.Context,
	vehicleID string,
	plateNumber string,
	vehicleType string,
	capacity int32,
) error {

	s.log.Infow(
		"service update vehicle request received",
		"vehicle_id", vehicleID,
	)

	if err := validateVehicleMetadata(plateNumber, capacity); err != nil {

		s.log.Warnw(
			"vehicle metadata validation failed",
			"vehicle_id", vehicleID,
			"error", err,
		)

		return err
	}

	err := s.repo.UpdateVehicle(
		ctx,
		vehicleID,
		plateNumber,
		vehicleType,
		capacity,
	)

	if err != nil {

		s.log.Errorw(
			"failed to update vehicle",
			"vehicle_id", vehicleID,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"vehicle updated successfully",
		"vehicle_id", vehicleID,
	)

	return nil
}

func (s *FleetService) UpdateVehicleStatus(
	ctx context.Context,
	vehicleID string,
	status string,
) error {

	s.log.Infow(
		"service update vehicle status request received",
		"vehicle_id", vehicleID,
		"status", status,
	)

	if err := validateStatus(status); err != nil {

		s.log.Warnw(
			"invalid vehicle status provided",
			"vehicle_id", vehicleID,
			"status", status,
			"error", err,
		)

		return err
	}

	err := s.repo.UpdateVehicleStatus(ctx, vehicleID, status)
	if err != nil {

		s.log.Errorw(
			"failed to update vehicle status",
			"vehicle_id", vehicleID,
			"status", status,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"vehicle status updated successfully",
		"vehicle_id", vehicleID,
		"status", status,
	)

	return nil
}

func (s *FleetService) UpdateDriver(
	ctx context.Context,
	driverID string,
	firstName string,
	lastName string,
) error {

	s.log.Infow(
		"service update driver request received",
		"driver_id", driverID,
	)

	if err := validateDriverName(firstName, lastName); err != nil {

		s.log.Warnw(
			"driver name validation failed",
			"driver_id", driverID,
			"error", err,
		)

		return err
	}

	err := s.repo.UpdateDriver(
		ctx,
		driverID,
		firstName,
		lastName,
	)

	if err != nil {

		s.log.Errorw(
			"failed to update driver",
			"driver_id", driverID,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"driver updated successfully",
		"driver_id", driverID,
	)

	return nil
}

func (s *FleetService) UpdateDriverStatus(
	ctx context.Context,
	driverID string,
	status string,
) error {

	s.log.Infow(
		"service update driver status request received",
		"driver_id", driverID,
		"status", status,
	)

	if err := validateStatus(status); err != nil {

		s.log.Warnw(
			"invalid driver status provided",
			"driver_id", driverID,
			"status", status,
			"error", err,
		)

		return err
	}

	err := s.repo.UpdateDriverStatus(ctx, driverID, status)
	if err != nil {

		s.log.Errorw(
			"failed to update driver status",
			"driver_id", driverID,
			"status", status,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"driver status updated successfully",
		"driver_id", driverID,
		"status", status,
	)

	return nil
}

func (s *FleetService) StartTrip(
	ctx context.Context,
	tripID string,
) error {

	s.log.Infow(
		"service start trip request received",
		"trip_id", tripID,
	)

	trip, err := s.repo.GetTrip(ctx, tripID)
	if err != nil {

		s.log.Errorw(
			"failed to fetch trip",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	if err := validateTripCanStart(ctx, s.repo, trip); err != nil {

		s.log.Warnw(
			"trip start validation failed",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	err = s.repo.StartTrip(ctx, tripID)
	if err != nil {

		s.log.Errorw(
			"failed to start trip",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"trip started successfully",
		"trip_id", tripID,
	)

	return nil
}

func (s *FleetService) CompleteTrip(
	ctx context.Context,
	tripID string,
) error {

	s.log.Infow(
		"service complete trip request received",
		"trip_id", tripID,
	)

	trip, err := s.repo.GetTrip(ctx, tripID)
	if err != nil {

		s.log.Errorw(
			"failed to fetch trip",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	if err := validateTripCanComplete(trip); err != nil {

		s.log.Warnw(
			"trip completion validation failed",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	err = s.repo.CompleteTrip(ctx, tripID)
	if err != nil {

		s.log.Errorw(
			"failed to complete trip",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"trip completed successfully",
		"trip_id", tripID,
	)

	return nil
}

// Supporter methods
func ensureVehicleAvailable(ctx context.Context, repo domain.FleetRepository, vehicleID string) error {

	active, err := repo.HasActiveTripForVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	if active {
		return errors.New("vehicle already has an active trip")
	}

	return nil
}

func ensureDriverAvailable(ctx context.Context, repo domain.FleetRepository, driverID string) error {

	active, err := repo.HasActiveTripForDriver(ctx, driverID)
	if err != nil {
		return err
	}

	if active {
		return errors.New("driver already has an active trip")
	}

	return nil
}

func validateTripOwnership(
	ctx context.Context,
	repo domain.FleetRepository,
	orgID string,
	vehicleID string,
	driverID string,
) error {

	v, err := repo.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	if v.OrganizationID != orgID {
		return errors.New("vehicle does not belong to organization")
	}

	d, err := repo.GetDriver(ctx, driverID)
	if err != nil {
		return err
	}

	if d.OrganizationID != orgID {
		return errors.New("driver does not belong to organization")
	}

	return nil
}

func validateStatus(status string) error {

	switch status {

	case "active", "inactive":
		return nil

	default:
		return errors.New("invalid status value")
	}
}

func validateVehicleMetadata(
	plateNumber string,
	capacity int32,
) error {

	if plateNumber == "" {
		return errors.New("plate number cannot be empty")
	}

	if capacity < 0 {
		return errors.New("capacity cannot be negative")
	}

	return nil
}

func validateDriverName(firstName, lastName string) error {

	if firstName == "" && lastName == "" {
		return errors.New("driver name cannot be empty")
	}

	return nil
}

func validateTripCanStart(
	ctx context.Context,
	repo domain.FleetRepository,
	trip *domain.Trip,
) error {

	if trip.Status != "created" {
		return errors.New("trip cannot be started from current state")
	}

	vehicle, err := repo.GetVehicle(ctx, trip.VehicleID)
	if err != nil {
		return err
	}

	if vehicle.Status != "active" {
		return errors.New("vehicle is not active")
	}

	driver, err := repo.GetDriver(ctx, trip.DriverID)
	if err != nil {
		return err
	}

	if driver.Status != "active" {
		return errors.New("driver is not active")
	}

	activeVehicleTrip, err := repo.HasActiveTripForVehicle(ctx, trip.VehicleID)
	if err != nil {
		return err
	}

	if activeVehicleTrip {
		return errors.New("vehicle already has an active trip")
	}

	activeDriverTrip, err := repo.HasActiveTripForDriver(ctx, trip.DriverID)
	if err != nil {
		return err
	}

	if activeDriverTrip {
		return errors.New("driver already has an active trip")
	}

	return nil
}

func validateTripCanComplete(
	trip *domain.Trip,
) error {

	if trip.Status != "active" {
		return errors.New("trip must be active to complete")
	}

	return nil
}
