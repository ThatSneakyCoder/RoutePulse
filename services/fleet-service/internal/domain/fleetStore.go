package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrNotFound          = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type FleetStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewFleetStore(db *sql.DB, log *zap.SugaredLogger) *FleetStore {
	return &FleetStore{
		db:  db,
		log: log,
	}
}

func (s *FleetStore) CreateVehicle(ctx context.Context, v *Vehicle) error {

	s.log.Infow(
		"creating vehicle",
		"vehicle_id", v.ID,
		"organization_id", v.OrganizationID,
		"plate_number", v.PlateNumber,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO vehicles (
		vehicle_id,
		organization_id,
		plate_number,
		vehicle_type,
		capacity,
		status
	)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		v.ID,
		v.OrganizationID,
		v.PlateNumber,
		v.VehicleType,
		v.Capacity,
		v.Status,
	).Scan(&v.CreatedAt)

	if err != nil {
		s.log.Errorw(
			"failed to create vehicle",
			"error", err,
			"vehicle_id", v.ID,
			"organization_id", v.OrganizationID,
		)
		return err
	}

	s.log.Infow(
		"vehicle created successfully",
		"vehicle_id", v.ID,
	)

	return nil
}

func (s *FleetStore) ListVehicles(ctx context.Context, organizationID string) ([]*Vehicle, error) {

	s.log.Infow(
		"listing vehicles",
		"organization_id", organizationID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		vehicle_id,
		organization_id,
		plate_number,
		vehicle_type,
		capacity,
		status,
		created_at
	FROM vehicles
	WHERE organization_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		s.log.Errorw(
			"failed to query vehicles",
			"error", err,
			"organization_id", organizationID,
		)
		return nil, err
	}

	defer rows.Close()

	vehicles := []*Vehicle{}

	for rows.Next() {

		var v Vehicle

		err := rows.Scan(
			&v.ID,
			&v.OrganizationID,
			&v.PlateNumber,
			&v.VehicleType,
			&v.Capacity,
			&v.Status,
			&v.CreatedAt,
		)

		if err != nil {
			s.log.Errorw(
				"failed to scan vehicle row",
				"error", err,
			)
			return nil, err
		}

		vehicles = append(vehicles, &v)
	}

	if err = rows.Err(); err != nil {
		s.log.Errorw(
			"row iteration error while listing vehicles",
			"error", err,
		)
		return nil, err
	}

	s.log.Infow(
		"vehicles fetched successfully",
		"organization_id", organizationID,
		"count", len(vehicles),
	)

	return vehicles, nil
}

func (s *FleetStore) ListDrivers(ctx context.Context, orgID string) ([]*Driver, error) {

	s.log.Infow(
		"listing drivers",
		"organization_id", orgID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		driver_id,
		organization_id,
		first_name,
		last_name,
		vehicle_id,
		status,
		created_at
	FROM drivers
	WHERE organization_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {

		s.log.Errorw(
			"failed to query drivers",
			"organization_id", orgID,
			"error", err,
		)

		return nil, err
	}

	defer rows.Close()

	drivers := []*Driver{}

	for rows.Next() {

		var d Driver

		err := rows.Scan(
			&d.ID,
			&d.OrganizationID,
			&d.FirstName,
			&d.LastName,
			&d.VehicleID,
			&d.Status,
			&d.CreatedAt,
		)

		if err != nil {

			s.log.Errorw(
				"failed to scan driver row",
				"error", err,
			)

			return nil, err
		}

		drivers = append(drivers, &d)
	}

	if err = rows.Err(); err != nil {

		s.log.Errorw(
			"row iteration error while listing drivers",
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

func (s *FleetStore) CreateDriver(ctx context.Context, d *Driver) error {

	s.log.Infow(
		"creating driver",
		"driver_id", d.ID,
		"organization_id", d.OrganizationID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO drivers (
		driver_id,
		organization_id,
		first_name,
		last_name,
		vehicle_id,
		status
	)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		d.ID,
		d.OrganizationID,
		d.FirstName,
		d.LastName,
		d.VehicleID,
		d.Status,
	).Scan(&d.CreatedAt)

	if err != nil {

		s.log.Errorw(
			"failed to create driver",
			"driver_id", d.ID,
			"organization_id", d.OrganizationID,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"driver created successfully",
		"driver_id", d.ID,
	)

	return nil
}

func (s *FleetStore) GetVehicle(ctx context.Context, vehicleID string) (*Vehicle, error) {

	s.log.Infow(
		"fetching vehicle",
		"vehicle_id", vehicleID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		vehicle_id,
		organization_id,
		plate_number,
		vehicle_type,
		capacity,
		status,
		created_at
	FROM vehicles
	WHERE vehicle_id = $1
	`

	var v Vehicle

	err := s.db.QueryRowContext(ctx, query, vehicleID).Scan(
		&v.ID,
		&v.OrganizationID,
		&v.PlateNumber,
		&v.VehicleType,
		&v.Capacity,
		&v.Status,
		&v.CreatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			s.log.Warnw(
				"vehicle not found",
				"vehicle_id", vehicleID,
			)
			return nil, ErrNotFound
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

	return &v, nil
}

func (s *FleetStore) CreateTrip(ctx context.Context, t *Trip) error {

	s.log.Infow(
		"creating trip",
		"trip_id", t.ID,
		"organization_id", t.OrganizationID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO trips (
		trip_id,
		organization_id,
		vehicle_id,
		driver_id,
		status
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		t.ID,
		t.OrganizationID,
		t.VehicleID,
		t.DriverID,
		t.Status,
	).Scan(&t.CreatedAt)

	if err != nil {

		s.log.Errorw(
			"failed to create trip",
			"trip_id", t.ID,
			"error", err,
		)

		return err
	}

	s.log.Infow(
		"trip created successfully",
		"trip_id", t.ID,
	)

	return nil
}

func (s *FleetStore) ListTrips(ctx context.Context, orgID string) ([]*Trip, error) {

	s.log.Infow(
		"listing trips",
		"organization_id", orgID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		trip_id,
		organization_id,
		vehicle_id,
		driver_id,
		status,
		start_time,
		end_time,
		created_at
	FROM trips
	WHERE organization_id=$1
	`

	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {

		s.log.Errorw(
			"failed to query trips",
			"error", err,
		)

		return nil, err
	}

	defer rows.Close()

	var trips []*Trip

	for rows.Next() {

		var t Trip

		err := rows.Scan(
			&t.ID,
			&t.OrganizationID,
			&t.VehicleID,
			&t.DriverID,
			&t.Status,
			&t.StartTime,
			&t.EndTime,
			&t.CreatedAt,
		)

		if err != nil {

			s.log.Errorw(
				"failed to scan trip",
				"error", err,
			)

			return nil, err
		}

		trips = append(trips, &t)
	}

	if err = rows.Err(); err != nil {

		s.log.Errorw(
			"trip rows iteration error",
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"trips fetched successfully",
		"organization_id", orgID,
		"count", len(trips),
	)

	return trips, nil
}

func (s *FleetStore) HasActiveTripForVehicle(ctx context.Context, vehicleID string) (bool, error) {

	s.log.Infow(
		"checking active trip for vehicle",
		"vehicle_id", vehicleID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM trips
		WHERE vehicle_id = $1
		AND status = 'active'
	)
	`

	var exists bool

	err := s.db.QueryRowContext(ctx, query, vehicleID).Scan(&exists)
	if err != nil {

		s.log.Errorw(
			"failed to check vehicle active trip",
			"vehicle_id", vehicleID,
			"error", err,
		)

		return false, err
	}

	s.log.Infow(
		"vehicle active trip check complete",
		"vehicle_id", vehicleID,
		"has_active_trip", exists,
	)

	return exists, nil
}

func (s *FleetStore) HasActiveTripForDriver(ctx context.Context, driverID string) (bool, error) {

	s.log.Infow(
		"checking active trip for driver",
		"driver_id", driverID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM trips
		WHERE driver_id = $1
		AND status = 'active'
	)
	`

	var exists bool

	err := s.db.QueryRowContext(ctx, query, driverID).Scan(&exists)
	if err != nil {

		s.log.Errorw(
			"failed to check driver active trip",
			"driver_id", driverID,
			"error", err,
		)

		return false, err
	}

	s.log.Infow(
		"driver active trip check complete",
		"driver_id", driverID,
		"has_active_trip", exists,
	)

	return exists, nil
}

func (s *FleetStore) GetDriver(ctx context.Context, driverID string) (*Driver, error) {

	s.log.Infow(
		"fetching driver",
		"driver_id", driverID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		driver_id,
		organization_id,
		first_name,
		last_name,
		vehicle_id,
		status,
		created_at
	FROM drivers
	WHERE driver_id = $1
	`

	var d Driver

	err := s.db.QueryRowContext(ctx, query, driverID).Scan(
		&d.ID,
		&d.OrganizationID,
		&d.FirstName,
		&d.LastName,
		&d.VehicleID,
		&d.Status,
		&d.CreatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {

			s.log.Warnw(
				"driver not found",
				"driver_id", driverID,
			)

			return nil, ErrNotFound
		}

		s.log.Errorw(
			"failed to fetch driver",
			"driver_id", driverID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"driver fetched successfully",
		"driver_id", driverID,
	)

	return &d, nil
}

func (s *FleetStore) UpdateVehicle(
	ctx context.Context,
	vehicleID string,
	plateNumber string,
	vehicleType string,
	capacity int32,
) error {

	s.log.Infow(
		"updating vehicle metadata",
		"vehicle_id", vehicleID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE vehicles
	SET
		plate_number = $2,
		vehicle_type = $3,
		capacity = $4
	WHERE vehicle_id = $1
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
		vehicleID,
		plateNumber,
		vehicleType,
		capacity,
	)

	if err != nil {

		s.log.Errorw(
			"failed to update vehicle metadata",
			"vehicle_id", vehicleID,
			"error", err,
		)

		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {

		s.log.Warnw(
			"vehicle not found while updating metadata",
			"vehicle_id", vehicleID,
		)

		return ErrNotFound
	}

	s.log.Infow(
		"vehicle metadata updated successfully",
		"vehicle_id", vehicleID,
	)

	return nil
}

func (s *FleetStore) UpdateVehicleStatus(
	ctx context.Context,
	vehicleID string,
	status string,
) error {

	s.log.Infow(
		"updating vehicle status",
		"vehicle_id", vehicleID,
		"status", status,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE vehicles
	SET status = $2
	WHERE vehicle_id = $1
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
		vehicleID,
		status,
	)

	if err != nil {

		s.log.Errorw(
			"failed to update vehicle status",
			"vehicle_id", vehicleID,
			"status", status,
			"error", err,
		)

		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {

		s.log.Warnw(
			"vehicle not found while updating status",
			"vehicle_id", vehicleID,
		)

		return ErrNotFound
	}

	s.log.Infow(
		"vehicle status updated successfully",
		"vehicle_id", vehicleID,
		"status", status,
	)

	return nil
}

func (s *FleetStore) UpdateDriver(
	ctx context.Context,
	driverID string,
	firstName string,
	lastName string,
) error {

	s.log.Infow(
		"updating driver profile",
		"driver_id", driverID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE drivers
	SET
		first_name = $2,
		last_name = $3
	WHERE driver_id = $1
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
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

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {

		s.log.Warnw(
			"driver not found while updating profile",
			"driver_id", driverID,
		)

		return ErrNotFound
	}

	s.log.Infow(
		"driver updated successfully",
		"driver_id", driverID,
	)

	return nil
}

func (s *FleetStore) UpdateDriverStatus(
	ctx context.Context,
	driverID string,
	status string,
) error {

	s.log.Infow(
		"updating driver status",
		"driver_id", driverID,
		"status", status,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE drivers
	SET status = $2
	WHERE driver_id = $1
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
		driverID,
		status,
	)

	if err != nil {

		s.log.Errorw(
			"failed to update driver status",
			"driver_id", driverID,
			"status", status,
			"error", err,
		)

		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {

		s.log.Warnw(
			"driver not found while updating status",
			"driver_id", driverID,
		)

		return ErrNotFound
	}

	s.log.Infow(
		"driver status updated successfully",
		"driver_id", driverID,
		"status", status,
	)

	return nil
}

func (s *FleetStore) StartTrip(ctx context.Context, tripID string) error {

	s.log.Infow("starting trip", "trip_id", tripID)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE trips
	SET
		status = 'active',
		start_time = now()
	WHERE trip_id = $1
	AND status = 'created'
	`

	res, err := s.db.ExecContext(ctx, query, tripID)
	if err != nil {

		s.log.Errorw(
			"failed to start trip",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {

		s.log.Warnw(
			"trip not found or already started",
			"trip_id", tripID,
		)

		return ErrNotFound
	}

	s.log.Infow(
		"trip started successfully",
		"trip_id", tripID,
	)

	return nil
}

func (s *FleetStore) CompleteTrip(ctx context.Context, tripID string) error {

	s.log.Infow("completing trip", "trip_id", tripID)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE trips
	SET
		status = 'completed',
		end_time = now()
	WHERE trip_id = $1
	AND status = 'active'
	`

	res, err := s.db.ExecContext(ctx, query, tripID)
	if err != nil {

		s.log.Errorw(
			"failed to complete trip",
			"trip_id", tripID,
			"error", err,
		)

		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {

		s.log.Warnw(
			"trip not found or not active",
			"trip_id", tripID,
		)

		return ErrNotFound
	}

	s.log.Infow(
		"trip completed successfully",
		"trip_id", tripID,
	)

	return nil
}

func (s *FleetStore) GetTrip(
	ctx context.Context,
	tripID string,
) (*Trip, error) {

	s.log.Infow(
		"fetching trip",
		"trip_id", tripID,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		trip_id,
		organization_id,
		vehicle_id,
		driver_id,
		status,
		start_time,
		end_time,
		created_at
	FROM trips
	WHERE trip_id = $1
	`

	var t Trip

	err := s.db.QueryRowContext(ctx, query, tripID).Scan(
		&t.ID,
		&t.OrganizationID,
		&t.VehicleID,
		&t.DriverID,
		&t.Status,
		&t.StartTime,
		&t.EndTime,
		&t.CreatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {

			s.log.Warnw(
				"trip not found",
				"trip_id", tripID,
			)

			return nil, ErrNotFound
		}

		s.log.Errorw(
			"failed to fetch trip",
			"trip_id", tripID,
			"error", err,
		)

		return nil, err
	}

	s.log.Infow(
		"trip fetched successfully",
		"trip_id", tripID,
		"vehicle_id", t.VehicleID,
		"driver_id", t.DriverID,
		"status", t.Status,
	)

	return &t, nil
}
