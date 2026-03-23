package main

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
)

type config struct {
	env      string
	grpcAddr string
	db       struct {
		user     string
		password string
		name     string
		host     string
		port     string
		sslMode  string
	}
}

func loadConfig() config {
	var cfg config

	// set env values
	cfg.env = env.GetString("ENV", "development")

	// set grpc port number
	cfg.grpcAddr = env.GetString("TRACKING_SERVICE_PORT", "9093")

	cfg.db.user = env.GetString("POSTGRES_USER", "")
	cfg.db.password = env.GetString("POSTGRES_PASSWORD", "")
	cfg.db.name = env.GetString("POSTGRES_DB", "")
	cfg.db.host = env.GetString("POSTGRES_HOST", "postgres-tracking")
	cfg.db.port = env.GetString("POSTGRES_PORT", "5432")
	cfg.db.sslMode = env.GetString("POSTGRES_SSLMODE", "disable")

	return cfg
}
