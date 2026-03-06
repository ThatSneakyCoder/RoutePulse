package main

import (
	"time"

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
	jwtAuthConfig struct {
		tokenConfig struct {
			secret string
			exp    time.Duration
			iss    string
			aud    string
		}
	}

	brevoMailConfig struct {
		host      string
		port      int
		username  string
		apiKey    string
		fromEmail string
		exp       time.Time
	}
}

func loadConfig() config {
	var cfg config

	// set env values
	cfg.env = env.GetString("ENV", "development")

	// set grpc port number
	cfg.grpcAddr = env.GetString("IDENTITY_SERVICE_PORT", "9090")

	// set database values
	cfg.db.user = env.GetString("POSTGRES_USER", "")
	cfg.db.password = env.GetString("POSTGRES_PASSWORD", "")
	cfg.db.name = env.GetString("POSTGRES_DB", "")

	cfg.db.host = "postgres"
	cfg.db.port = "5432"
	cfg.db.sslMode = "disable"

	// set jwt values
	cfg.jwtAuthConfig.tokenConfig.secret = env.GetString("AUTH_TOKEN_SECRET", "sampele_auth_token_string")
	cfg.jwtAuthConfig.tokenConfig.exp = time.Hour * 24 * 3
	cfg.jwtAuthConfig.tokenConfig.iss = env.GetString("AUTH_TOKEN_ISS", "routepulse.identity")
	cfg.jwtAuthConfig.tokenConfig.aud = env.GetString("AUTH_TOKEN_AUD", "routepulse.api")

	// mail config
	cfg.brevoMailConfig.host = env.GetString("BREVO_SMTP_HOST", "smtp-relay.brevo.com")
	cfg.brevoMailConfig.port = env.GetInt("BREVO_SMTP_PORT", 587)
	cfg.brevoMailConfig.username = env.GetString("BREVO_SMTP_USERNAME", "")
	cfg.brevoMailConfig.apiKey = env.GetString("BREVO_SMTP_KEY", "")
	cfg.brevoMailConfig.fromEmail = env.GetString("", "shubh22.srivastava@gmail.com")

	return cfg
}
