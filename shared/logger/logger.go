package logger

import "go.uber.org/zap"

type Config struct {
	ServiceName string
	Env         string // development | production
}

func Init(cfg Config) *zap.SugaredLogger {
	var l *zap.Logger
	var err error

	if cfg.Env == "production" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}

	return l.With(
		zap.String("service", cfg.ServiceName),
	).Sugar()
}
