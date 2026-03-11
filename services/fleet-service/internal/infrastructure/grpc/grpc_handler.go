package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/services/fleet-service/internal/service"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/fleet"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedFleetServiceServer
	log     *zap.SugaredLogger
	service *service.FleetService
}

func NewGRPCHandler(server *grpc.Server, svc *service.FleetService, log *zap.SugaredLogger) {
	handler := &gRPCHandler{
		log:     log,
		service: svc,
	}

	pb.RegisterFleetServiceServer(server, handler)
}
