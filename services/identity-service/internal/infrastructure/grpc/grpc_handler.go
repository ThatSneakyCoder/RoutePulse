package grpc

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/service"
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedIdentityServiceServer
	log     *zap.SugaredLogger
	service *service.IdentityService
}

func NewGRPCHandler(server *grpc.Server, svc *service.IdentityService, log *zap.SugaredLogger) {
	handler := &gRPCHandler{
		service: svc,
		log:     log,
	}
	pb.RegisterIdentityServiceServer(server, handler)
}

func (h *gRPCHandler) RegisterUser(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {

	h.log.Infow("RegisterUser called",
		"email", req.Email,
	)

	user := &domain.User{
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Email:     req.Email,
		Password:  req.Password,
	}

	created, err := h.service.RegisterUser(ctx, user)
	if err != nil {
		h.log.Errorw("failed to register user", "err", err)
		return nil, err
	}

	h.log.Infow("user registered successfully", "user_id", created.ID)

	return &pb.RegisterUserResponse{
		User: UserToProto(created),
	}, nil
}

// func (h *gRPCHandler) GetUserByEmail(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
// 	h.log.Infow("GetUserByEmail called",
// 		"email", req.Email,
// 	)

// 	user := &domain.User{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	}

// 	retrievedUser, err := h.service.GetUserByEmail(ctx, user)
// 	if err != nil {
// 		h.log.Errorw("failed to get user by email", "err", err)
// 		return nil, err
// 	}

// 	h.log.Infow("obatined user successfully", "user_id", retrievedUser.ID)

// 	// return &pb.LoginResponse{
		
// 	// }

// 	return nil
// }
