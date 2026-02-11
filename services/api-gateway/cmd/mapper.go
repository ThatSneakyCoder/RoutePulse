package main

import (
	pb "github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	pbOrg "github.com/ThatSneakyCoder/RoutePulse/shared/proto/organization"
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

func (p *ForgotPasswordRequest) toProto() *pb.ForgotPasswordRequest {
	return &pb.ForgotPasswordRequest{
		Email: p.Email,
	}
}

func (p *ResetPasswordRequest) toProto() *pb.ResetPasswordRequest {
	return &pb.ResetPasswordRequest{
		Email: p.Email,
		Token: p.Token,
		NewPassword: p.NewPassword,
	}
}

func (p *CreateOrganizationRequest) toProto(userID string) *pbOrg.CreateOrganizationRequest {
	return &pbOrg.CreateOrganizationRequest{
		Name: p.Name,
		OwnerUserId: userID,
	}
}

func (p *ListUserOrganizationsRequest) toProto(userID string) *pbOrg.ListUserOrganizationsRequest {
	return &pbOrg.ListUserOrganizationsRequest{
		UserId: userID,
	}
}