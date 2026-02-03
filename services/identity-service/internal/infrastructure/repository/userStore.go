package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
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

func (s *UserStore) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	const query = `
		INSERT INTO users (
			id,
			email,
			password_hash,
			first_name,
			last_name,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			email,
			first_name,
			last_name,
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
		user.Password,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

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
