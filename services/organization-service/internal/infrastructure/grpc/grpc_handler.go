package grpc

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/organization-service/internal/service"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/organization"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedOrganizationServiceServer
	log     *zap.SugaredLogger
	service *service.OrganizationService
}

func NewGRPCHandler(server *grpc.Server, svc *service.OrganizationService, log *zap.SugaredLogger) {
	handler := &gRPCHandler{
		log:     log,
		service: svc,
	}

	pb.RegisterOrganizationServiceServer(server, handler)
}

func (h *gRPCHandler) CreateOrganization(
	ctx context.Context,
	req *pb.CreateOrganizationRequest,
) (*pb.CreateOrganizationResponse, error) {

	h.log.Infow("CreateOrganization called",
		"owner_user_id", req.OwnerUserId,
		"name", req.Name,
	)

	org, err := h.service.CreateOrganization(
		ctx,
		req.Name,
		req.OwnerUserId,
	)
	if err != nil {
		h.log.Errorw("failed to create organization",
			"owner_user_id", req.OwnerUserId,
			"name", req.Name,
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("organization created successfully",
		"organization_id", org.ID,
		"owner_user_id", org.OwnerUserID,
	)

	return &pb.CreateOrganizationResponse{
		Organization: OrganizationToProto(org),
	}, nil
}

func (h *gRPCHandler) ListUserOrganizations(
	ctx context.Context,
	req *pb.ListUserOrganizationsRequest,
) (*pb.ListUserOrganizationsResponse, error) {

	h.log.Infow("ListUserOrganizations called",
		"user_id", req.UserId,
	)

	orgs, err := h.service.ListOrganizationsByUserID(ctx, req.UserId)
	if err != nil {
		h.log.Errorw("failed to list user organizations",
			"user_id", req.UserId,
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("user organizations fetched successfully",
		"user_id", req.UserId,
		"count", len(orgs),
	)

	return &pb.ListUserOrganizationsResponse{
		Organizations: OrganizationsToProto(orgs),
	}, nil
}
