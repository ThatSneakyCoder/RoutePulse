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

func (h *gRPCHandler) GetOrganization(
	ctx context.Context,
	req *pb.GetOrganizationRequest,
) (*pb.GetOrganizationResponse, error) {

	h.log.Infow("GetOrganization called",
		"organization_id", req.OrganizationId,
	)

	org, err := h.service.GetOrganization(ctx, req.OrganizationId)
	if err != nil {
		h.log.Errorw("failed to fetch organization",
			"organization_id", req.OrganizationId,
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("organization fetched successfully",
		"organization_id", org.ID,
	)

	return &pb.GetOrganizationResponse{
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

func (h *gRPCHandler) ListOrganizationMembers(
	ctx context.Context,
	req *pb.ListOrganizationMembersRequest,
) (*pb.ListOrganizationMembersResponse, error) {
	h.log.Infow("ListOrganizationMembers called",
		"organization_id", req.OrganizationId,
	)

	members, err := h.service.ListOrganizationMembers(ctx, req.OrganizationId)
	if err != nil {
		h.log.Errorw("failed to list organization members",
			"organization_id", req.OrganizationId,
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("organization members fetched successfully",
		"organization_id", req.OrganizationId,
		"count", len(members),
	)

	return &pb.ListOrganizationMembersResponse{
		Members: OrganizationMembersToProto(members),
	}, nil
}

func (h *gRPCHandler) AddMember(
	ctx context.Context,
	req *pb.AddMemberRequest,
) (*pb.AddMemberResponse, error) {

	h.log.Infow("AddMember called",
		"organization_id", req.OrganizationId,
		"user_id", req.UserId,
		"role", req.Role,
	)

	err := h.service.AddMember(
		ctx,
		req.OrganizationId,
		req.UserId,
		req.Role,
	)

	if err != nil {
		h.log.Errorw("failed to add member",
			"organization_id", req.OrganizationId,
			"user_id", req.UserId,
			"err", err,
		)
		return nil, err
	}

	return &pb.AddMemberResponse{
		Success: true,
	}, nil
}

func (h *gRPCHandler) GetTotalMembers(
	ctx context.Context,
	req *pb.GetTotalMembersRequest,
) (*pb.GetTotalMembersResponse, error) {

	h.log.Infow("GetTotalMembers called",
		"owner_user_id", req.OwnerUserId,
	)

	count, err := h.service.GetTotalMembers(ctx, req.OwnerUserId)
	if err != nil {
		h.log.Errorw("failed to get total members",
			"owner_user_id", req.OwnerUserId,
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("total members fetched successfully",
		"owner_user_id", req.OwnerUserId,
		"count", count,
	)

	return &pb.GetTotalMembersResponse{
		Count: count,
	}, nil
}

func (h *gRPCHandler) RemoveMember(
	ctx context.Context,
	req *pb.RemoveMemberRequest,
) (*pb.RemoveMemberResponse, error) {

	h.log.Infow("RemoveMember called",
		"organization_id", req.OrganizationId,
		"user_id", req.UserId,
	)

	err := h.service.RemoveMember(
		ctx,
		req.OrganizationId,
		req.UserId,
	)

	if err != nil {
		h.log.Errorw("failed to remove member",
			"organization_id", req.OrganizationId,
			"user_id", req.UserId,
			"err", err,
		)
		return nil, err
	}

	return &pb.RemoveMemberResponse{
		Success: true,
	}, nil
}

func (h *gRPCHandler) UpdateMemberRole(
	ctx context.Context,
	req *pb.UpdateMemberRoleRequest,
) (*pb.UpdateMemberRoleResponse, error) {

	h.log.Infow("UpdateMemberRole called",
		"organization_id", req.OrganizationId,
		"user_id", req.UserId,
		"role", req.Role,
	)

	err := h.service.UpdateMemberRole(
		ctx,
		req.OrganizationId,
		req.UserId,
		req.Role,
	)

	if err != nil {
		h.log.Errorw("failed to update member role",
			"organization_id", req.OrganizationId,
			"user_id", req.UserId,
			"err", err,
		)
		return nil, err
	}

	return &pb.UpdateMemberRoleResponse{
		Success: true,
	}, nil
}
