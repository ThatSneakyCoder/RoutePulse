package main

import (
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
)

func (p *CreateUserRequest) toProto() *pb.RegisterUserRequest {
	return &pb.RegisterUserRequest{
		Firstname: p.FirstName,
		Lastname:  p.LastName,
		Email:     p.Email,
		Password:  p.Password,
	}
}

func (p *LoginUserRequest) toProto() *pb.LoginRequest {
	return &pb.LoginRequest{
		Email:    p.Email,
		Password: p.Password,
	}
}

func (p *ValidateTokenRequest) toProto() *pb.ValidateTokenRequest {
	return &pb.ValidateTokenRequest{
		AccessToken: p.AccessToken,
	}
}

func (p *VerifyUserEmailRequest) toProto() *pb.VerifyUserEmailRequest {
	return &pb.VerifyUserEmailRequest{
		Email: p.Email,
		Token: p.Token,
	}
}
