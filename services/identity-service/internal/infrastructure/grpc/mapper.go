package grpc

import (
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/auth"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(user *domain.User) *pb.User {
	return &pb.User{
		UserId:    user.ID.String(),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func LoginResponseToProto(token *auth.Token) *pb.LoginResponse {
	return &pb.LoginResponse{
		AccessToken: token.AccessToken,
		ExpiresAt:   token.ExpiresAt.Unix(), // TODO: should be timestamp - or something similar - need to figure this out
	}
}
