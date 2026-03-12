package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrNotFound          = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type UserStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewUserStore(db *sql.DB, log *zap.SugaredLogger) *UserStore {
	return &UserStore{
		db:  db,
		log: log,
	}
}

func (s *UserStore) Create(ctx context.Context, user *User) (*User, error) {
	const query = `
			INSERT INTO users (
			id,
			email,
			first_name,
			last_name,
			password_hash,
			is_active,
			is_verified,
			verify_email_token_hash,
			verify_email_token_expires_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING
			id,
			email,
			first_name,
			last_name,
			is_active,
			is_verified,
			verify_email_token_hash,
			verify_email_token_expires_at,
			created_at,
			updated_at
	`

	s.log.Debugw("creating user in database",
		"user_id", user.ID,
		"email", user.Email,
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password.hash,
		user.IsActive,
		user.IsVerified,
		user.EmailVerification.TokenHash,
		user.EmailVerification.ExpiresAt,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsVerified,
		&user.EmailVerification.TokenHash,
		&user.EmailVerification.ExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// we will not rely on in-memory struct because, database is source of truth. Hence, we will re-hydrate EmailVerification with returned values from query

	// Note: I didn't follow what I wrote above because, I don't understand how to write proper code for the above

	if err != nil {
		translated := translatePostgresError(err)

		switch translated {
		case ErrDuplicateEmail:
			s.log.Warnw("duplicate email constraint violation",
				"email", user.Email,
			)
		case ErrDuplicateUsername:
			s.log.Warnw("duplicate username constraint violation")
		default:
			s.log.Errorw("failed to create user",
				"user_id", user.ID,
				"email", user.Email,
				"err", err,
			)
		}

		return nil, translated
	}

	s.log.Infow("user persisted successfully",
		"user_id", user.ID,
	)

	return user, nil
}

func (s *UserStore) DeleteUserByEmail(ctx context.Context, email string) error {
	s.log.Debugw("deleting user by email", "email", email)

	const query = `
		delete from users
		where
			email = $1
		and
			is_verified = false
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, email)
	if err != nil {
		s.log.Debugw("user could not be deleted", "email", email)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {

	s.log.Debugw("fetching user by email", "email", email)

	const query = `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			is_active,
			is_verified,
			verify_email_token_hash,
			verify_email_token_expires_at,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	var hash []byte
	var tokenHash *string
	var tokenExpiresAt *time.Time

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&hash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsVerified,
		&tokenHash,
		&tokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Debugw("user not found", "email", email)
			return nil, ErrNotFound
		}

		s.log.Errorw("failed to fetch user by email",
			"email", email,
			"err", err,
		)
		return nil, err
	}

	user.Password.SetHash(hash)

	if tokenHash != nil || tokenExpiresAt != nil {
		user.EmailVerification = &EmailVerification{
			TokenHash: tokenHash,
			ExpiresAt: tokenExpiresAt,
		}
	}

	s.log.Debugw("user fetched successfully",
		"user_id", user.ID,
	)

	return &user, nil
}

func (s *UserStore) MarkEmailVerified(ctx context.Context, user *User) error {
	const query = `
	UPDATE users
	SET
		is_verified = true,
		verify_email_token_hash = NULL,
		verify_email_token_expires_at = NULL,
		updated_at = $1
	WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

func (s *UserStore) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {

	s.log.Debugw("fetching user by id", "id", id.String())

	const query = `
		SELECT
			id,
			email,
			first_name,
			last_name,
			is_active,
			is_verified,
			verify_email_token_hash,
			verify_email_token_expires_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	var tokenHash *string
	var tokenExpiresAt *time.Time

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsVerified,
		&tokenHash,
		&tokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Debugw("user not found", "user_id", id)
			return nil, ErrNotFound
		}

		s.log.Errorw("db error fetching user",
			"user_id", id,
			"err", err,
		)
		return nil, err
	}

	if tokenHash != nil || tokenExpiresAt != nil {
		user.EmailVerification = &EmailVerification{
			TokenHash: tokenHash,
			ExpiresAt: tokenExpiresAt,
		}
	}

	s.log.Debugw("user fetched successfully",
		"user_id", user.ID,
	)

	return &user, nil
}

func (s *UserStore) Upsert(
	ctx context.Context,
	reset *PasswordReset,
) error {

	query := `
		INSERT INTO password_resets (
			id, user_id, token_hash, expires_at
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_at = NOW()
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		reset.ID,
		reset.UserID,
		reset.TokenHash,
		reset.ExpiresAt,
	)
	return err
}

func (s *UserStore) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*PasswordReset, error) {

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM password_resets
		WHERE user_id = $1
	`

	var reset PasswordReset
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&reset.ID,
		&reset.UserID,
		&reset.TokenHash,
		&reset.ExpiresAt,
		&reset.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &reset, nil
}

func (s *UserStore) DeleteByUserID(
	ctx context.Context,
	userID uuid.UUID,
) error {
	query := `
		DELETE FROM password_resets WHERE user_id = $1
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		userID,
	)

	return err
}

func (s *UserStore) UpdatePassword(
	ctx context.Context,
	user *User,
) error {

	query := `
		UPDATE users
		SET password_hash = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		user.Password.hash,
		user.ID,
	)

	return err
}

func (s *UserStore) GetUsersByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]*User, error) {

	s.log.Debugw("fetching users by ids",
		"count", len(ids),
	)

	const query = `
		SELECT
			id,
			email,
			first_name,
			last_name,
			is_active,
			is_verified,
			created_at,
			updated_at
		FROM users
		WHERE id = ANY($1)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		s.log.Errorw("failed to query users by ids",
			"err", err,
		)
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {

		var u User

		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.IsActive,
			&u.IsVerified,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			s.log.Errorw("failed to scan user row",
				"err", err,
			)
			return nil, err
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		s.log.Errorw("user rows iteration failed",
			"err", err,
		)
		return nil, err
	}

	s.log.Debugw("users fetched successfully",
		"count", len(users),
	)

	return users, nil
}

func translatePostgresError(err error) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}

	switch pqErr.Constraint {
	case "users_email_key":
		return ErrDuplicateEmail
	case "users_username_key":
		return ErrDuplicateUsername
	default:
		return err
	}
}
