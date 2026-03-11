package domain

import (
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
