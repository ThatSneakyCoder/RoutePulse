package main

import (
	pbFleet "github.com/ThatSneakyCoder/RoutePulse/shared/proto/fleet"
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
		Email:       p.Email,
		Token:       p.Token,
		NewPassword: p.NewPassword,
	}
}

func (p *CreateOrganizationRequest) toProto(userID string) *pbOrg.CreateOrganizationRequest {
	return &pbOrg.CreateOrganizationRequest{
		Name:        p.Name,
		OwnerUserId: userID,
	}
}

func (p *ListUserOrganizationsRequest) toProto(userID string) *pbOrg.ListUserOrganizationsRequest {
	return &pbOrg.ListUserOrganizationsRequest{
		UserId: userID,
	}
}

func (p *GetOrganizationMembersListRequest) toProto() *pbOrg.ListOrganizationMembersRequest {
	return &pbOrg.ListOrganizationMembersRequest{
		OrganizationId: p.OrganizationID,
	}
}

func (p *UserIDsRequest) toProto() *pb.GetUsersByIDsRequest {
	return &pb.GetUsersByIDsRequest{
		UserIds: p.UserIDs,
	}
}

func (p *GetOrganizationRequest) toProto() *pbOrg.GetOrganizationRequest {
	return &pbOrg.GetOrganizationRequest{
		OrganizationId: p.OrganizationID,
	}
}

func (p *GetUserByEmailRequest) toProto() *pb.GetUserByEmailRequest {
	return &pb.GetUserByEmailRequest{
		Email: p.Email,
	}
}

func (p *AddOrganizationMemberRequest) toProto() *pbOrg.AddMemberRequest {
	return &pbOrg.AddMemberRequest{
		OrganizationId: p.OrganizationID,
		UserId:         p.UserID,
		Role:           p.Role,
	}
}

func (p *RemoveOrganizationMemberRequest) toProto() *pbOrg.RemoveMemberRequest {
	return &pbOrg.RemoveMemberRequest{
		OrganizationId: p.OrganizationID,
		UserId:         p.UserID,
	}
}

func (p *UpdateOrganizationMemberRoleRequest) toProto() *pbOrg.UpdateMemberRoleRequest {
	return &pbOrg.UpdateMemberRoleRequest{
		OrganizationId: p.OrganizationID,
		UserId:         p.UserID,
		Role:           p.Role,
	}
}

func (p *CreateVehicleRequest) toProto() *pbFleet.CreateVehicleRequest {
	return &pbFleet.CreateVehicleRequest{
		OrganizationId: p.OrganizationID,
		PlateNumber:    p.PlateNumber,
		VehicleType:    p.VehicleType,
		Capacity:       p.Capacity,
	}
}

func (p *ListVehiclesRequest) toProto() *pbFleet.ListVehiclesRequest {
	return &pbFleet.ListVehiclesRequest{
		OrganizationId: p.OrganizationID,
	}
}

func (p *GetVehicleRequest) toProto() *pbFleet.GetVehicleRequest {
	return &pbFleet.GetVehicleRequest{
		VehicleId: p.VehicleID,
	}
}

func (p *CreateDriverRequest) toProto() *pbFleet.CreateDriverRequest {
	return &pbFleet.CreateDriverRequest{
		OrganizationId: p.OrganizationID,
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		VehicleId:      p.VehicleID,
	}
}

func (p *ListDriversRequest) toProto() *pbFleet.ListDriversRequest {
	return &pbFleet.ListDriversRequest{
		OrganizationId: p.OrganizationID,
	}
}

func (p *CreateTripRequest) toProto() *pbFleet.CreateTripRequest {
	return &pbFleet.CreateTripRequest{
		OrganizationId: p.OrganizationID,
		VehicleId:      p.VehicleID,
		DriverId:       p.DriverID,
	}
}

func (p *ListTripsRequest) toProto() *pbFleet.ListTripsRequest {
	return &pbFleet.ListTripsRequest{
		OrganizationId: p.OrganizationID,
	}
}

func (p *UpdateVehicleRequest) toProto() *pbFleet.UpdateVehicleRequest {
	return &pbFleet.UpdateVehicleRequest{
		VehicleId:   p.VehicleID,
		PlateNumber: p.PlateNumber,
		VehicleType: p.VehicleType,
		Capacity:    p.Capacity,
	}
}

func (p *UpdateVehicleStatusRequest) toProto() *pbFleet.UpdateVehicleStatusRequest {
	return &pbFleet.UpdateVehicleStatusRequest{
		VehicleId: p.VehicleID,
		Status:    p.Status,
	}
}

func (p *UpdateDriverRequest) toProto() *pbFleet.UpdateDriverRequest {
	return &pbFleet.UpdateDriverRequest{
		DriverId:  p.DriverID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
}

func (p *UpdateDriverStatusRequest) toProto() *pbFleet.UpdateDriverStatusRequest {
	return &pbFleet.UpdateDriverStatusRequest{
		DriverId: p.DriverID,
		Status:   p.Status,
	}
}

func (p *StartTripRequest) toProto() *pbFleet.StartTripRequest {
	return &pbFleet.StartTripRequest{
		TripId: p.TripID,
	}
}

func (p *CompleteTripRequest) toProto() *pbFleet.CompleteTripRequest {
	return &pbFleet.CompleteTripRequest{
		TripId: p.TripID,
	}
}

func toGetTotalMembersGRPC(req GetTotalMembersRequest) *pbOrg.GetTotalMembersRequest {
	return &pbOrg.GetTotalMembersRequest{
		OwnerUserId: req.OwnerUserID,
	}
}