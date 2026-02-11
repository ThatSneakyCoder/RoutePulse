package main

import (
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/api-gateway/internal/infrastructure/grpc"
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	"github.com/ThatSneakyCoder/RoutePulse/shared/logger"
)

//	@title			Fleet Management System API
//	@version		1.0
//	@description	REST API documentation for Fleet Management System
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath		/v1
//	@schemes		http https

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Enter the token with the `Bearer ` prefix, e.g. `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),
	}

	// Logger
	log := logger.Init(logger.Config{
		ServiceName: "api-gateway",
		Env:         cfg.env,
	})
	defer log.Sync()

	log.Infow("starting api gateway",
		"addr", cfg.addr,
		"env", cfg.env,
	)

	// Initalize grpc client
	// identity service
	identityClient, err := grpc.NewIdentityServiceClient()
	if err != nil {
		log.Fatalw(
			"failed to connect to identity service",
			"service", "identity-service",
			"err", err,
		)
	}
	defer identityClient.Close()

	log.Infow(
		"connected to downstream service",
		"service", "identity-service",
	)

	// identity service
	analyticsClient, err := grpc.NewAnalyticsServiceClient()
	log.Infow(
		"connected to downstream service",
		"service", "analytics-service",
	)
	defer analyticsClient.Close()

	// organization service
	organizationClient, err := grpc.NewOrganizationServiceClient()
	log.Infow(
		"connected to downstream service",
		"service", "organization-service",
	)
	defer organizationClient.Close()

	// Metrics
	metrics := newMetrics()

	app := &application{
		config:             cfg,
		log:                log,
		metrics:            metrics,
		identityClient:     identityClient,
		analyticsClient:    analyticsClient,
		organizationClient: organizationClient,
	}

	app.limiters = rateLimiters{
		global: newLimiterEntry(rateLimiterConfig{
			RequestsPerTimeFrame: 300,
			TimeFrame:            time.Minute,
			Enabled:              true,
		}),
		register: newLimiterEntry(rateLimiterConfig{
			RequestsPerTimeFrame: 5,
			TimeFrame:            time.Minute,
			Enabled:              true,
		}),
		login: newLimiterEntry(rateLimiterConfig{
			RequestsPerTimeFrame: 10,
			TimeFrame:            time.Minute,
			Enabled:              true,
		}),
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatalw("api gateway stopped with error", "err", err)
	}

	log.Info("api gateway stopped gracefully")
}
