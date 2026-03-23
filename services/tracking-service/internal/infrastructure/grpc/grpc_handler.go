package grpc

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/tracking-service/internal/service"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/tracking"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
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

func (h *gRPCHandler) GetTripCurrentLocation(
	ctx context.Context,
	req *pb.GetTripCurrentLocationRequest,
) (*pb.GetTripCurrentLocationResponse, error) {

	h.log.Infow(
		"grpc get trip current location request received",
		"trip_id", req.TripId,
	)

	location, err := h.service.GetTripCurrentLocation(ctx, req.TripId)
	if err != nil {
		h.log.Errorw(
			"grpc get trip current location failed",
			"trip_id", req.TripId,
			"error", err,
		)

		if err == domain.ErrNotFound {
			return nil, grpcstatus.Error(codes.NotFound, err.Error())
		}

		return nil, grpcstatus.Error(codes.Internal, err.Error())
	}

	return &pb.GetTripCurrentLocationResponse{
		Location: &pb.TripCurrentLocation{
			TripId:     location.TripID,
			DriverId:   location.DriverID,
			VehicleId:  location.VehicleID,
			Latitude:   location.Latitude,
			Longitude:  location.Longitude,
			RecordedAt: location.RecordedAt.Unix(),
			Connection: location.Connection,
		},
	}, nil
}

func (h *gRPCHandler) GetTripLocationHistory(
	ctx context.Context,
	req *pb.GetTripLocationHistoryRequest,
) (*pb.GetTripLocationHistoryResponse, error) {

	h.log.Infow(
		"grpc get trip location history request received",
		"trip_id", req.TripId,
		"limit", req.Limit,
	)

	points, err := h.service.GetTripLocationHistory(ctx, req.TripId, req.Limit)
	if err != nil {
		h.log.Errorw(
			"grpc get trip location history failed",
			"trip_id", req.TripId,
			"error", err,
		)

		if err == domain.ErrNotFound {
			return nil, grpcstatus.Error(codes.NotFound, err.Error())
		}

		return nil, grpcstatus.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.TripLocationHistoryPoint, 0, len(points))
	for _, point := range points {
		result = append(result, &pb.TripLocationHistoryPoint{
			Latitude:   point.Latitude,
			Longitude:  point.Longitude,
			RecordedAt: point.RecordedAt.Unix(),
		})
	}

	return &pb.GetTripLocationHistoryResponse{
		TripId: req.TripId,
		Points: result,
	}, nil
}

func (h *gRPCHandler) GetTripGeometry(
	ctx context.Context,
	req *pb.GetTripGeometryRequest,
) (*pb.GetTripGeometryResponse, error) {

	h.log.Infow(
		"grpc get trip geometry request received",
		"trip_id", req.TripId,
	)

	geometry, err := h.service.GetTripGeometry(ctx, req.TripId)
	if err != nil {
		h.log.Errorw(
			"grpc get trip geometry failed",
			"trip_id", req.TripId,
			"error", err,
		)

		if err == domain.ErrNotFound {
			return nil, grpcstatus.Error(codes.NotFound, err.Error())
		}

		return nil, grpcstatus.Error(codes.Internal, err.Error())
	}

	return &pb.GetTripGeometryResponse{
		Geometry: &pb.TripGeometry{
			TripId:          geometry.TripID,
			PlannedGeometry: toProtoCoordinates(geometry.PlannedGeometry),
			ActualGeometry:  toProtoCoordinates(geometry.ActualGeometry),
		},
	}, nil
}

func toProtoCoordinates(points []domain.Coordinate) []*pb.Coordinate {
	result := make([]*pb.Coordinate, 0, len(points))

	for _, point := range points {
		result = append(result, &pb.Coordinate{
			Latitude:  point.Latitude,
			Longitude: point.Longitude,
		})
	}

	return result
}
