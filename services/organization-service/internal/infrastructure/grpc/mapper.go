package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/services/organization-service/internal/domain"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/organization"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrganizationToProto(org *domain.Organization) *pb.Organization {
	return &pb.Organization{
		OrganizationId: org.ID.String(),
		Name:           org.Name,
		OwnerUserId:    org.OwnerUserID.String(),
		IsActive:       org.IsActive,
		CreatedAt:      timestamppb.New(org.CreatedAt),
		UpdatedAt:      timestamppb.New(org.UpdatedAt),
	}
}

func OrganizationsToProto(orgs []*domain.Organization) []*pb.Organization {
	result := make([]*pb.Organization, 0, len(orgs))
	for _, org := range orgs {
		result = append(result, OrganizationToProto(org))
	}
	return result
}
