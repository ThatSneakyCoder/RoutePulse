package grpc

import (
	"context"

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

func (h *gRPCHandler) CreateVehicle(
	ctx context.Context,
	req *pb.CreateVehicleRequest,
) (*pb.CreateVehicleResponse, error) {

	h.log.Infow(
		"grpc CreateVehicle request received",
		"organization_id", req.OrganizationId,
		"plate_number", req.PlateNumber,
	)

	v, err := h.service.CreateVehicle(
		ctx,
		req.OrganizationId,
		req.PlateNumber,
		req.VehicleType,
		req.Capacity,
	)

	if err != nil {
		h.log.Errorw(
			"grpc CreateVehicle failed",
			"error", err,
			"organization_id", req.OrganizationId,
		)
		return nil, err
	}

	h.log.Infow(
		"grpc CreateVehicle success",
		"vehicle_id", v.ID,
	)

	return &pb.CreateVehicleResponse{
		Vehicle: &pb.Vehicle{
			VehicleId:      v.ID,
			OrganizationId: v.OrganizationID,
			PlateNumber:    v.PlateNumber,
			VehicleType:    v.VehicleType,
			Capacity:       v.Capacity,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt.Unix(),
		},
	}, nil
}

func (h *gRPCHandler) ListVehicles(
	ctx context.Context,
	req *pb.ListVehiclesRequest,
) (*pb.ListVehiclesResponse, error) {

	h.log.Infow(
		"grpc ListVehicles request received",
		"organization_id", req.OrganizationId,
	)

	vehicles, err := h.service.ListVehicles(ctx, req.OrganizationId)
	if err != nil {

		h.log.Errorw(
			"grpc ListVehicles failed",
			"error", err,
			"organization_id", req.OrganizationId,
		)

		return nil, err
	}

	res := make([]*pb.Vehicle, 0, len(vehicles))

	for _, v := range vehicles {

		res = append(res, &pb.Vehicle{
			VehicleId:      v.ID,
			OrganizationId: v.OrganizationID,
			PlateNumber:    v.PlateNumber,
			VehicleType:    v.VehicleType,
			Capacity:       v.Capacity,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt.Unix(),
		})
	}

	h.log.Infow(
		"grpc ListVehicles success",
		"organization_id", req.OrganizationId,
		"count", len(res),
	)

	return &pb.ListVehiclesResponse{
		Vehicles: res,
	}, nil
}
