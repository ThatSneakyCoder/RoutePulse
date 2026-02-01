package main

import (
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
	"go.uber.org/zap"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
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
