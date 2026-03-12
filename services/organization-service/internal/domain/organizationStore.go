package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound          = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type OrganizationStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewOrganizationStore(db *sql.DB, log *zap.SugaredLogger) *OrganizationStore {
	return &OrganizationStore{
		db:  db,
		log: log,
	}
}

func (s *OrganizationStore) Create(
	ctx context.Context,
	org *Organization,
) (*Organization, error) {

	const insertOrgQuery = `
		INSERT INTO organizations (
			id,
			name,
			owner_user_id,
			is_active,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			name,
			owner_user_id,
			is_active,
			created_at,
			updated_at
	`

	s.log.Debugw("creating organization in database",
		"organization_id", org.ID,
		"owner_user_id", org.OwnerUserID,
		"name", org.Name,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		insertOrgQuery,
		org.ID,
		org.Name,
		org.OwnerUserID,
		org.IsActive,
		org.CreatedAt,
		org.UpdatedAt,
	).Scan(
		&org.ID,
		&org.Name,
		&org.OwnerUserID,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		s.log.Errorw("failed to create organization",
			"organization_id", org.ID,
			"owner_user_id", org.OwnerUserID,
			"err", err,
		)
		return nil, err
	}

	// owner membership insert (standalone helper)
	if err := insertOwnerMember(ctx, s.db, org.ID, org.OwnerUserID); err != nil {
		s.log.Errorw("failed to insert organization owner membership",
			"organization_id", org.ID,
			"user_id", org.OwnerUserID,
			"err", err,
		)
		return nil, err
	}

	s.log.Infow("organization persisted successfully",
		"organization_id", org.ID,
	)

	return org, nil
}

func (s *OrganizationStore) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]*Organization, error) {

	s.log.Debugw("listing organizations for user",
		"user_id", userID,
	)

	const query = `
		SELECT
			o.id,
			o.name,
			o.owner_user_id,
			o.is_active,
			o.created_at,
			o.updated_at
		FROM organizations o
		JOIN organization_members om
			ON o.id = om.organization_id
		WHERE om.user_id = $1
		  AND o.is_active = true
		ORDER BY o.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		s.log.Errorw("failed to list organizations for user",
			"user_id", userID,
			"err", err,
		)
		return nil, err
	}
	defer rows.Close()

	orgs := []*Organization{}

	for rows.Next() {
		var org Organization

		if err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.OwnerUserID,
			&org.IsActive,
			&org.CreatedAt,
			&org.UpdatedAt,
		); err != nil {
			s.log.Errorw("failed to scan organization row",
				"user_id", userID,
				"err", err,
			)
			return nil, err
		}

		orgs = append(orgs, &org)
	}

	if err := rows.Err(); err != nil {
		s.log.Errorw("organization row iteration failed",
			"user_id", userID,
			"err", err,
		)
		return nil, err
	}

	s.log.Debugw("organizations fetched successfully",
		"user_id", userID,
		"count", len(orgs),
	)

	return orgs, nil
}

func insertOwnerMember(
	ctx context.Context,
	db *sql.DB,
	orgID uuid.UUID,
	userID uuid.UUID,
) error {

	const query = `
		INSERT INTO organization_members (
			organization_id,
			user_id,
			role,
			joined_at
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.ExecContext(
		ctx,
		query,
		orgID,
		userID,
		"OWNER",
		time.Now().UTC(),
	)

	return err
}

func (s *OrganizationStore) ListMembersByOrganizationID(
	ctx context.Context,
	orgID uuid.UUID,
) ([]*OrganizationMember, error) {

	s.log.Debugw("listing organization members from database",
		"organization_id", orgID,
	)

	const query = `
		SELECT
			organization_id,
			user_id,
			role,
			joined_at
		FROM organization_members
		WHERE organization_id = $1
		ORDER BY joined_at ASC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {
		s.log.Errorw("failed to query organization members",
			"organization_id", orgID,
			"err", err,
		)
		return nil, err
	}
	defer rows.Close()

	members := []*OrganizationMember{}

	for rows.Next() {
		var m OrganizationMember

		if err := rows.Scan(
			&m.OrganizationID,
			&m.UserID,
			&m.Role,
			&m.JoinedAt,
		); err != nil {
			s.log.Errorw("failed to scan organization member row",
				"organization_id", orgID,
				"err", err,
			)
			return nil, err
		}

		members = append(members, &m)
	}

	if err := rows.Err(); err != nil {
		s.log.Errorw("organization members row iteration failed",
			"organization_id", orgID,
			"err", err,
		)
		return nil, err
	}

	s.log.Debugw("organization members fetched successfully",
		"organization_id", orgID,
		"count", len(members),
	)

	return members, nil
}

func (s *OrganizationStore) GetByID(
	ctx context.Context,
	orgID uuid.UUID,
) (*Organization, error) {

	s.log.Debugw("fetching organization by id",
		"organization_id", orgID,
	)

	const query = `
		SELECT
			id,
			name,
			owner_user_id,
			is_active,
			created_at,
			updated_at
		FROM organizations
		WHERE id = $1
		  AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var org Organization

	err := s.db.QueryRowContext(ctx, query, orgID).Scan(
		&org.ID,
		&org.Name,
		&org.OwnerUserID,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warnw("organization not found",
				"organization_id", orgID,
			)
			return nil, ErrNotFound
		}

		s.log.Errorw("failed to fetch organization",
			"organization_id", orgID,
			"err", err,
		)
		return nil, err
	}

	s.log.Debugw("organization fetched successfully",
		"organization_id", org.ID,
	)

	return &org, nil
}

func (s *OrganizationStore) AddMember(
	ctx context.Context,
	orgID uuid.UUID,
	userID uuid.UUID,
	role string,
) error {

	const query = `
	INSERT INTO organization_members
	(
		organization_id,
		user_id,
		role,
		joined_at
	)
	VALUES ($1,$2,$3,$4)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		orgID,
		userID,
		role,
		time.Now().UTC(),
	)

	if err != nil {
		s.log.Errorw("failed to insert organization member",
			"organization_id", orgID,
			"user_id", userID,
			"err", err,
		)
		return err
	}

	s.log.Infow("organization member added",
		"organization_id", orgID,
		"user_id", userID,
	)

	return nil
}

func (s *OrganizationStore) RemoveMember(
	ctx context.Context,
	orgID uuid.UUID,
	userID uuid.UUID,
) error {

	const query = `
		DELETE FROM organization_members
		WHERE organization_id = $1
		AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		orgID,
		userID,
	)

	if err != nil {
		s.log.Errorw("failed to delete organization member",
			"organization_id", orgID,
			"user_id", userID,
			"err", err,
		)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	s.log.Infow("organization member removed",
		"organization_id", orgID,
		"user_id", userID,
	)

	return nil
}

func (s *OrganizationStore) UpdateMemberRole(
	ctx context.Context,
	orgID uuid.UUID,
	userID uuid.UUID,
	role string,
) error {

	const query = `
		UPDATE organization_members
		SET role = $1
		WHERE organization_id = $2
		AND user_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		role,
		orgID,
		userID,
	)

	if err != nil {
		s.log.Errorw("failed to update member role",
			"organization_id", orgID,
			"user_id", userID,
			"err", err,
		)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	s.log.Infow("organization member role updated",
		"organization_id", orgID,
		"user_id", userID,
		"role", role,
	)

	return nil
}
