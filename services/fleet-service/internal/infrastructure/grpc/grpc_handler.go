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

func (h *gRPCHandler) GetVehicle(
	ctx context.Context,
	req *pb.GetVehicleRequest,
) (*pb.GetVehicleResponse, error) {

	h.log.Infow(
		"grpc GetVehicle request received",
		"vehicle_id", req.VehicleId,
	)

	v, err := h.service.GetVehicle(ctx, req.VehicleId)
	if err != nil {

		h.log.Errorw(
			"grpc GetVehicle failed",
			"vehicle_id", req.VehicleId,
			"error", err,
		)

		return nil, err
	}

	resp := &pb.GetVehicleResponse{
		Vehicle: &pb.Vehicle{
			VehicleId:      v.ID,
			OrganizationId: v.OrganizationID,
			PlateNumber:    v.PlateNumber,
			VehicleType:    v.VehicleType,
			Capacity:       v.Capacity,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt.Unix(),
		},
	}

	h.log.Infow(
		"grpc GetVehicle success",
		"vehicle_id", v.ID,
	)

	return resp, nil
}

func (h *gRPCHandler) CreateDriver(
	ctx context.Context,
	req *pb.CreateDriverRequest,
) (*pb.CreateDriverResponse, error) {

	h.log.Infow(
		"grpc CreateDriver request received",
		"organization_id", req.OrganizationId,
	)

	var vehicleID *string

	if req.VehicleId != "" {
		vehicleID = &req.VehicleId
	}

	driver, err := h.service.CreateDriver(
		ctx,
		req.OrganizationId,
		req.FirstName,
		req.LastName,
		vehicleID,
	)

	if err != nil {

		h.log.Errorw(
			"grpc CreateDriver failed",
			"organization_id", req.OrganizationId,
			"error", err,
		)

		return nil, err
	}

	resp := &pb.CreateDriverResponse{
		Driver: &pb.Driver{
			DriverId:       driver.ID,
			OrganizationId: driver.OrganizationID,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			VehicleId:      req.VehicleId,
			Status:         driver.Status,
			CreatedAt:      driver.CreatedAt.Unix(),
		},
	}

	h.log.Infow(
		"grpc CreateDriver success",
		"driver_id", driver.ID,
	)

	return resp, nil
}

func (h *gRPCHandler) ListDrivers(
	ctx context.Context,
	req *pb.ListDriversRequest,
) (*pb.ListDriversResponse, error) {

	h.log.Infow(
		"grpc ListDrivers request received",
		"organization_id", req.OrganizationId,
	)

	drivers, err := h.service.ListDrivers(ctx, req.OrganizationId)
	if err != nil {

		h.log.Errorw(
			"grpc ListDrivers failed",
			"organization_id", req.OrganizationId,
			"error", err,
		)

		return nil, err
	}

	result := []*pb.Driver{}

	for _, d := range drivers {

		var vehicleID string

		if d.VehicleID != nil {
			vehicleID = *d.VehicleID
		}

		result = append(result, &pb.Driver{
			DriverId:       d.ID,
			OrganizationId: d.OrganizationID,
			FirstName:      d.FirstName,
			LastName:       d.LastName,
			VehicleId:      vehicleID,
			Status:         d.Status,
			CreatedAt:      d.CreatedAt.Unix(),
		})
	}

	resp := &pb.ListDriversResponse{
		Drivers: result,
	}

	h.log.Infow(
		"grpc ListDrivers success",
		"organization_id", req.OrganizationId,
		"count", len(result),
	)

	return resp, nil
}

func (h *gRPCHandler) CreateTrip(
	ctx context.Context,
	req *pb.CreateTripRequest,
) (*pb.CreateTripResponse, error) {

	h.log.Infow(
		"grpc CreateTrip request",
		"organization_id", req.OrganizationId,
	)

	trip, err := h.service.CreateTrip(
		ctx,
		req.OrganizationId,
		req.VehicleId,
		req.DriverId,
	)

	if err != nil {

		h.log.Errorw(
			"grpc CreateTrip failed",
			"error", err,
		)

		return nil, err
	}

	return &pb.CreateTripResponse{
		Trip: &pb.Trip{
			TripId:         trip.ID,
			OrganizationId: trip.OrganizationID,
			VehicleId:      trip.VehicleID,
			DriverId:       trip.DriverID,
			Status:         trip.Status,
			CreatedAt:      trip.CreatedAt.Unix(),
		},
	}, nil
}

func (h *gRPCHandler) ListTrips(
	ctx context.Context,
	req *pb.ListTripsRequest,
) (*pb.ListTripsResponse, error) {

	h.log.Infow(
		"grpc ListTrips request",
		"organization_id", req.OrganizationId,
	)

	trips, err := h.service.ListTrips(ctx, req.OrganizationId)
	if err != nil {

		h.log.Errorw(
			"grpc ListTrips failed",
			"error", err,
		)

		return nil, err
	}

	resp := []*pb.Trip{}

	for _, t := range trips {

		resp = append(resp, &pb.Trip{
			TripId:         t.ID,
			OrganizationId: t.OrganizationID,
			VehicleId:      t.VehicleID,
			DriverId:       t.DriverID,
			Status:         t.Status,
			CreatedAt:      t.CreatedAt.Unix(),
		})
	}

	return &pb.ListTripsResponse{
		Trips: resp,
	}, nil
}

func (h *gRPCHandler) UpdateVehicle(
	ctx context.Context,
	req *pb.UpdateVehicleRequest,
) (*pb.UpdateVehicleResponse, error) {

	h.log.Infow(
		"grpc update vehicle request received",
		"vehicle_id", req.VehicleId,
	)

	err := h.service.UpdateVehicle(
		ctx,
		req.VehicleId,
		req.PlateNumber,
		req.VehicleType,
		req.Capacity,
	)
	if err != nil {

		h.log.Errorw(
			"grpc update vehicle failed",
			"vehicle_id", req.VehicleId,
			"error", err,
		)

		return nil, err
	}

	h.log.Infow(
		"vehicle updated successfully",
		"vehicle_id", req.VehicleId,
	)

	return &pb.UpdateVehicleResponse{
		Success: true,
	}, nil
}

func (h *gRPCHandler) UpdateVehicleStatus(
	ctx context.Context,
	req *pb.UpdateVehicleStatusRequest,
) (*pb.UpdateVehicleStatusResponse, error) {

	h.log.Infow(
		"grpc update vehicle status request received",
		"vehicle_id", req.VehicleId,
		"status", req.Status,
	)

	err := h.service.UpdateVehicleStatus(ctx, req.VehicleId, req.Status)
	if err != nil {

		h.log.Errorw(
			"grpc update vehicle status failed",
			"vehicle_id", req.VehicleId,
			"status", req.Status,
			"error", err,
		)

		return nil, err
	}

	h.log.Infow(
		"vehicle status updated successfully",
		"vehicle_id", req.VehicleId,
		"status", req.Status,
	)

	return &pb.UpdateVehicleStatusResponse{
		Success: true,
	}, nil
}

func (h *gRPCHandler) UpdateDriver(
	ctx context.Context,
	req *pb.UpdateDriverRequest,
) (*pb.UpdateDriverResponse, error) {

	h.log.Infow(
		"grpc update driver request received",
		"driver_id", req.DriverId,
	)

	err := h.service.UpdateDriver(
		ctx,
		req.DriverId,
		req.FirstName,
		req.LastName,
	)

	if err != nil {

		h.log.Errorw(
			"grpc update driver failed",
			"driver_id", req.DriverId,
			"error", err,
		)

		return nil, err
	}

	h.log.Infow(
		"driver updated successfully",
		"driver_id", req.DriverId,
	)

	return &pb.UpdateDriverResponse{
		Success: true,
	}, nil
}

func (h *gRPCHandler) UpdateDriverStatus(
	ctx context.Context,
	req *pb.UpdateDriverStatusRequest,
) (*pb.UpdateDriverStatusResponse, error) {

	h.log.Infow(
		"grpc update driver status request received",
		"driver_id", req.DriverId,
		"status", req.Status,
	)

	err := h.service.UpdateDriverStatus(ctx, req.DriverId, req.Status)
	if err != nil {

		h.log.Errorw(
			"grpc update driver status failed",
			"driver_id", req.DriverId,
			"status", req.Status,
			"error", err,
		)

		return nil, err
	}

	h.log.Infow(
		"driver status updated successfully",
		"driver_id", req.DriverId,
		"status", req.Status,
	)

	return &pb.UpdateDriverStatusResponse{
		Success: true,
	}, nil
}

func (h *gRPCHandler) StartTrip(
	ctx context.Context,
	req *pb.StartTripRequest,
) (*pb.StartTripResponse, error) {

	h.log.Infow(
		"grpc start trip request received",
		"trip_id", req.TripId,
	)

	err := h.service.StartTrip(ctx, req.TripId)
	if err != nil {

		h.log.Errorw(
			"grpc start trip failed",
			"trip_id", req.TripId,
			"error", err,
		)

		return nil, err
	}

	h.log.Infow(
		"trip started successfully",
		"trip_id", req.TripId,
	)

	return &pb.StartTripResponse{Success: true}, nil
}

func (h *gRPCHandler) CompleteTrip(
	ctx context.Context,
	req *pb.CompleteTripRequest,
) (*pb.CompleteTripResponse, error) {

	h.log.Infow(
		"grpc complete trip request received",
		"trip_id", req.TripId,
	)

	err := h.service.CompleteTrip(ctx, req.TripId)
	if err != nil {

		h.log.Errorw(
			"grpc complete trip failed",
			"trip_id", req.TripId,
			"error", err,
		)

		return nil, err
	}

	h.log.Infow(
		"trip completed successfully",
		"trip_id", req.TripId,
	)

	return &pb.CompleteTripResponse{Success: true}, nil
}