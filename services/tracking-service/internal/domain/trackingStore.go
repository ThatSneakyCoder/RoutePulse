package domain

import (
	"database/sql"

	"go.uber.org/zap"
)

type TrackingStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewTrackingStore(db *sql.DB, log *zap.SugaredLogger) *TrackingStore {
	return &TrackingStore{
		db:  db,
		log: log,
	}
}
