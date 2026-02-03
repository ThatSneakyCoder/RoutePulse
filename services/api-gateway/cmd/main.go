package main

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	"go.uber.org/zap"
)

// @title           Fleet Management System API
// @version         1.0
// @description     REST API documentation for Fleet Management System
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath        /v1
// @schemes         http https

// @securityDefinitions.apiKey ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Info("Zap logger setup successfully")

	metrics := newMetrics()

	app := &application{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
	}

	mux := app.mount()
	logger.Info(app.run(mux))
}
