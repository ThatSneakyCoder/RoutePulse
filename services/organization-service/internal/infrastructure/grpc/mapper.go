package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/services/organization-service/internal/domain"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/organization"
)

func OrganizationToProto(org *domain.Organization) *pb.Organization {
	return &pb.Organization{
		OrganizationId: org.ID.String(),
		Name:           org.Name,
		OwnerUserId:    org.OwnerUserID.String(),
		IsActive:       org.IsActive,
		CreatedAt:      org.CreatedAt.Unix(),
		UpdatedAt:      org.UpdatedAt.Unix(),
	}
}

func OrganizationsToProto(orgs []*domain.Organization) []*pb.Organization {
	result := make([]*pb.Organization, 0, len(orgs))
	for _, org := range orgs {
		result = append(result, OrganizationToProto(org))
	}
	return result
}

func OrganizationMemberToProto(m *domain.OrganizationMember) *pb.OrganizationMember {
	return &pb.OrganizationMember{
		UserId:   m.UserID.String(),
		Role:     m.Role,
		JoinedAt: m.JoinedAt.Unix(),
	}
}

func OrganizationMembersToProto(members []*domain.OrganizationMember) []*pb.OrganizationMember {
	result := make([]*pb.OrganizationMember, 0, len(members))

	for _, m := range members {
		result = append(result, OrganizationMemberToProto(m))
	}

	return result
}
