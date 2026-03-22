package domain

import "time"

type AnalyticsEvent struct {
	EventTime  time.Time
	OrgID      string
	UserID     string
	EventType  string
	EntityType string
	EntityID   *string
	Service    string
	Route      string
	StatusCode uint16
	LatencyMs  uint32
	RequestID  string
}

type Event struct {
	EventType string
	UserID    string
	OrgID     string
	Service   string
	EventTime time.Time
}
