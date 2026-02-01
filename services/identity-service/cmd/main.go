package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/grpc"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/repository"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/service"
	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9090"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		// TODO: need to change logging with zap package
	}

	userRepo := repository.NewUserStore()
	identitySvc := service.NewIdentityService(userRepo)
	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, identitySvc)

	log.Printf("Starting grpc server identity service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
