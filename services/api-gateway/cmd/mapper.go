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

// func (p *LoginUserRequest) toProto() *pb.LoginRequest {
// 	return &pb.
// }
