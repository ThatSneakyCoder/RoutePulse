package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/auth"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/db"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/grpc"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/service"
	"github.com/ThatSneakyCoder/RoutePulse/shared/logger"
	grpcserver "google.golang.org/grpc"
)

func main() {
	cfg := loadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Logger
	log := logger.Init(logger.Config{
		ServiceName: "identity-service",
		Env:         cfg.env,
	})
	defer log.Sync()

	log.Infow("identity service starting",
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

	lis, err := net.Listen("tcp", ":"+cfg.grpcAddr)
	if err != nil {
		log.Fatalw("failed to listen",
			"addr", cfg.grpcAddr,
			"err", err,
		)
	}

	// Database
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
		log.Fatalw("failed to connect to database", "err", err)
	}

	log.Info("database connection established")

	JWTAuthenticator := auth.NewJWTAuthenticator(
		cfg.jwtAuthConfig.tokenConfig.secret,
		cfg.jwtAuthConfig.tokenConfig.aud,
		cfg.jwtAuthConfig.tokenConfig.iss,
		cfg.jwtAuthConfig.tokenConfig.exp,
	)

	// Wiring
	userRepo := domain.NewUserStore(dbConn, log)
	identitySvc := service.NewIdentityService(userRepo, log, JWTAuthenticator)

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, identitySvc, log)

	log.Infow("starting gRPC server",
		"addr", lis.Addr().String(),
	)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Errorw("gRPC server stopped with error", "err", err)
			cancel()
		}
	}()

	<-ctx.Done()

	log.Info("gracefully shutting down gRPC server")
	grpcServer.GracefulStop()
}
