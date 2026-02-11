package domain

import "time"

type AnalyticsEvent struct {
	EventTime  time.Time
	Service    string
	EventType  string
	UserID     string
	RequestID  string
	Route      string
	StatusCode uint16
	LatencyMs  uint32
}
