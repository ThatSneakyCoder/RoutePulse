package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/infrastructure/db"
	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/infrastructure/events"
	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/infrastructure/grpc"
	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/service"
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
		ServiceName: "analytics-service",
		Env:         cfg.env,
	})
	defer log.Sync()

	log.Infow("analytics service starting",
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

	// ClickHouse connection
	chConn, err := db.Connect(ctx, cfg.db.host, cfg.db.port, cfg.db.name, cfg.db.user, cfg.db.password)
	if err != nil {
		log.Fatalw("failed to connect to clickhouse", "err", err)
	}
	defer chConn.Close()

	log.Info("connected to clickhouse")

	// gRPC listener
	lis, err := net.Listen("tcp", ":"+cfg.grpcAddr)
	if err != nil {
		log.Fatalw("failed to listen",
			"addr", cfg.grpcAddr,
			"err", err,
		)
	}

	eventStore := domain.NewEventsStore(chConn, log)
	analyticsSvc := service.NewAnalyticsService(eventStore, log)

	// RabbitMQ connection
	rabbitmq, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	// Consume rabbitMQ events from other services
	identityConsumer := events.NewIdentityConsumer(rabbitmq, analyticsSvc, log)
	go identityConsumer.Listen()

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, log, analyticsSvc)

	log.Infow("starting analytics gRPC server",
		"addr", lis.Addr().String(),
	)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Errorw("gRPC server stopped with error", "err", err)
			cancel()
		}
	}()

	<-ctx.Done()

	log.Info("gracefully shutting down analytics gRPC server")
	grpcServer.GracefulStop()
}
