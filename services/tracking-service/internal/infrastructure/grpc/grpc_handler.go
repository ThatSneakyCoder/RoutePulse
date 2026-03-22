package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/tracking"
)

type gRPCHandler struct {
	pb.UnimplementedTrackingServiceServer
	log     *zap.SugaredLogger
	service *service.TrackingService
}

func NewGRPCHandler(server *grpc.Server, svc *service.TrackingService, log *zap.SugaredLogger) {
	handler := &gRPCHandler{
		log:     log,
		service: svc,
	}

	pb.RegisterTrackingServiceServer(server, handler)
}
