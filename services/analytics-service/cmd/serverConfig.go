package main

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
)

type config struct {
	env      string
	grpcAddr string
	db       struct {
		host     string
		port     string
		name     string
		user     string
		password string
	}
}

func loadConfig() config {
	var cfg config

	cfg.env = env.GetString("ENV", "development")
	cfg.grpcAddr = env.GetString("ANALYTICS_SERVICE_PORT", "9096")

	cfg.db.host = env.GetString("CLICKHOUSE_HOST", "clickhouse")
	cfg.db.port = env.GetString("CLICKHOUSE_PORT", "9000")
	cfg.db.name = env.GetString("CLICKHOUSE_DB", "analytics")
	cfg.db.password = env.GetString("CLICKHOUSE_PASSWORD", "guest")
	cfg.db.user = env.GetString("CLICKHOUSE_USER", "guest")

	return cfg
}
