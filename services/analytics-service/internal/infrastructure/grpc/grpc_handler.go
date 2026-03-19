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

// GetTotalMembers
func (h *gRPCHandler) GetTotalMembers(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.TotalMembersResponse, error) {

	h.log.Infow("GetTotalMembers gRPC request received")

	count, err := h.service.GetTotalMembers(ctx)
	if err != nil {
		h.log.Errorw("GetTotalMembers failed", "err", err)
		return nil, err
	}

	h.log.Infow("GetTotalMembers completed successfully", "count", count)

	return &pb.TotalMembersResponse{
		Count: count,
	}, nil
}

// GetActiveUsersToday
func (h *gRPCHandler) GetActiveUsersToday(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.ActiveUsersTodayResponse, error) {

	h.log.Infow("GetActiveUsersToday gRPC request received")

	count, err := h.service.GetActiveUsersToday(ctx)
	if err != nil {
		h.log.Errorw("GetActiveUsersToday failed", "err", err)
		return nil, err
	}

	h.log.Infow("GetActiveUsersToday completed successfully", "count", count)

	return &pb.ActiveUsersTodayResponse{
		Count: count,
	}, nil
}

func (h *gRPCHandler) GetRecentActivity(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.RecentActivityResponse, error) {

	h.log.Infow("GetRecentActivity gRPC request received")

	events, err := h.service.GetRecentActivity(ctx)
	if err != nil {
		h.log.Errorw("GetRecentActivity failed", "err", err)
		return nil, err
	}

	pbEvents := make([]*pb.Event, len(events))

	for i, e := range events {
		pbEvents[i] = &pb.Event{
			EventType: e.EventType,
			UserId:    e.UserID,
			OrgId:     e.OrgID,
			Service:   e.Service,
			EventTime: e.EventTime,
		}
	}

	h.log.Infow("GetRecentActivity completed successfully",
		"events_count", len(pbEvents),
	)

	return &pb.RecentActivityResponse{
		Events: pbEvents,
	}, nil
}
