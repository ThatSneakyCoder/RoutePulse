package grpc

import (
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
)

func UserToProto(user *domain.User) *pb.User {
	return &pb.User{
		UserId:    user.ID.String(),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Unix(),
	}
}