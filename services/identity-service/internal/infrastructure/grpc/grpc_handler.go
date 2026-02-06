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
	}

	created, err := h.service.RegisterUser(ctx, user, req.Password)
	if err != nil {
		h.log.Errorw("failed to register user", "err", err)
		return nil, err
	}

	h.log.Infow("user registered successfully", "user_id", created.ID)

	return &pb.RegisterUserResponse{
		User: UserToProto(created),
	}, nil
}

func (h *gRPCHandler) VerifyUserEmail(ctx context.Context, req *pb.VerifyUserEmailRequest) (*pb.VerifyUserEmailResponse, error) {
	h.log.Infow("VerifyUserEmail called",
		"email", req.Email,
	)

	userID, verified, err := h.service.VerifyUserEmail(
		ctx,
		req.Email,
		req.Token,
	)
	if err != nil {
		h.log.Warnw("VerifyUserEmail failed",
			"email", req.Email,
			"err", err,
		)
		return nil, err
	}

	return &pb.VerifyUserEmailResponse{
		UserId:     userID,
		IsVerified: verified,
	}, nil
}

func (h *gRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	h.log.Infow("GetUserByEmail called",
		"email", req.Email,
	)

	user, token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.log.Warnw("login failed",
			"email", req.Email,
			"err", err,
		)
		return nil, err // already a gRPC status error
	}

	h.log.Infow("login successful",
		"user_id", user.ID,
	)
	return LoginResponseToProto(token), nil
}

func (h *gRPCHandler) ValidateToken(
	ctx context.Context,
	req *pb.ValidateTokenRequest,
) (*pb.ValidateTokenResponse, error) {

	h.log.Infow("ValidateTokenAndGetUser called")

	user, err := h.service.ValidateJWTTokenAndGetUser(ctx, req.AccessToken)
	if err != nil {
		h.log.Errorw("token validation failed",
			"err", err,
		)
		return nil, err
	}

	h.log.Infow("token validated successfully",
		"user_id", user.ID,
	)

	return &pb.ValidateTokenResponse{
		Valid: true,
		User:  UserToProto(user),
	}, nil
}
