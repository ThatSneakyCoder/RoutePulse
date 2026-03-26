package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/infrastructure/db"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/infrastructure/events"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/infrastructure/grpc"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/service"
	"github.com/ThatSneakyCoder/RoutePulse/shared/logger"
	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"

	grpcserver "google.golang.org/grpc"
)

func main() {
	cfg := loadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Logger
	log := logger.Init(logger.Config{
		ServiceName: "tracking-service",
		Env:         cfg.env,
	})
	defer log.Sync()

	log.Infow("tracking service starting",
		"env", cfg.env,
		"grpc_addr", cfg.grpcAddr,
	)

	// Handle shutdown signals
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Info("shutdown signal received")
		cancel()
	}()

	dsn :=
		"postgres://" +
			cfg.db.user + ":" +
			cfg.db.password +
			"@" + cfg.db.host + ":" + cfg.db.port + "/" +
			cfg.db.name +
			"?sslmode=" + cfg.db.sslMode

	dbConn, err := db.New(db.Config{
		DSN:          dsn,
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  5 * time.Minute,
	})
	if err != nil {
		log.Fatalw("failed to connect to tracking database", "err", err)
	}
	defer dbConn.Close()

	log.Info("database connection established in tracking service")

	// gRPC listener
	lis, err := net.Listen("tcp", ":"+cfg.grpcAddr)
	if err != nil {
		log.Fatalw("failed to listen",
			"addr", cfg.grpcAddr,
			"err", err,
		)
	}

	// RabbitMQ connection
	rabbitmq, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	trackingRepo := domain.NewTrackingStore(dbConn, log)

	trackingSvc := service.NewTrackingService(trackingRepo, log)
	trackingConsumer := events.NewTrackingConsumer(rabbitmq, trackingSvc, log)

	if err := trackingConsumer.Listen(); err != nil {
		log.Fatalw("failed to start tracking consumer", "err", err)
	}

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, trackingSvc, log)

	log.Infow("starting tracking gRPC server",
		"addr", lis.Addr().String(),
	)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Errorw("gRPC server stopped with error", "err", err)
			cancel()
		}
	}()

	<-ctx.Done()

	log.Info("gracefully shutting down tracking gRPC server")
	grpcServer.GracefulStop()
}
