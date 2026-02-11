package ratelimiter

import "time"

type Limiter interface {
	Allow(string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
