package grpc

import (
	"context"
	"fmt"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/service"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedIdentityServiceServer
	service *service.IdentityService
}

func NewGRPCHandler(server *grpc.Server, svc *service.IdentityService) {
	handler := &gRPCHandler{service: svc}
	pb.RegisterIdentityServiceServer(server, handler)
}

func (h *gRPCHandler) RegisterUser(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {

	// user := &domain.User{
	// 	Email: req.Email,
	// }

	// created, err := h.service.RegisterUser(ctx, user)
	// if err != nil {
	// 	return nil, err
	// }

	// return &pb.RegisterUserResponse{
	// 	// UserId: created.ID.String(),
	// 	// Email:  created.Email,
	// }, nil

	fmt.Println("HIT detected in identity microservice")

	return nil, nil
}
