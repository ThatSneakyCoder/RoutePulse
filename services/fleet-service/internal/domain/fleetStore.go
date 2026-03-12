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
