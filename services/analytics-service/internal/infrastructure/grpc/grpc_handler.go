package grpc

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/analytics-service/internal/service"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/analytics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type gRPCHandler struct {
	pb.UnimplementedAnalyticsServiceServer
	log     *zap.SugaredLogger
	service *service.AnalyticsService
}

func NewGRPCHandler(server *grpc.Server, log *zap.SugaredLogger, service *service.AnalyticsService) {
	handler := &gRPCHandler{
		log:     log,
		service: service,
	}
	pb.RegisterAnalyticsServiceServer(server, handler)
}

func (h *gRPCHandler) GetVehiclesInTransit(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.VehiclesInTransitResponse, error) {

	h.log.Infow("GetVehiclesInTransit gRPC request received")

	count, err := h.service.GetVehiclesInTransit(ctx)
	if err != nil {
		h.log.Errorw("GetVehiclesInTransit failed",
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("GetVehiclesInTransit completed successfully",
		"count", count,
	)

	return &pb.VehiclesInTransitResponse{
		Count: count,
	}, nil
}

func (h *gRPCHandler) GetTripsToday(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.TripsTodayResponse, error) {

	h.log.Infow("GetTripsToday gRPC request received")

	count, err := h.service.GetTripsToday(ctx)
	if err != nil {
		h.log.Errorw("GetTripsToday failed",
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("GetTripsToday completed successfully",
		"count", count,
	)

	return &pb.TripsTodayResponse{
		Count: count,
	}, nil
}
