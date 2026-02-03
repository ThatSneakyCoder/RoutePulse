package main

import "os"

type config struct {
	env      string
	grpcAddr string

	db struct {
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

	cfg.env = os.Getenv("ENV")
	if cfg.env == "" {
		cfg.env = "development"
	}

	cfg.grpcAddr = os.Getenv("IDENTITY_SERVICE_PORT")
	if cfg.grpcAddr == "" {
		cfg.grpcAddr = "9090"
	}

	cfg.db.user = os.Getenv("POSTGRES_USER")
	cfg.db.password = os.Getenv("POSTGRES_PASSWORD")
	cfg.db.name = os.Getenv("POSTGRES_DB")

	cfg.db.host = "postgres"
	cfg.db.port = "5432"
	cfg.db.sslMode = "disable"

	return cfg
}
