package main

import "time"

// this is for auth middleware to store authenticated user only. Not directly related to grpc
type AuthenticatedUser struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstname" validate:"required,max=255"`
	LastName  string `json:"lastname" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=3,max=72"`
}

// CreateUserResponse is used only for Swagger documentation
type CreateUserResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName string `json:"firstname" example:"John"`
	LastName  string `json:"lastname" example:"Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	CreatedAt int64  `json:"created_at" example:"1706950000"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type ValidateTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type ValidateTokenResponse struct {
	Valid bool         `json:"valid"`
	User  UserResponse `json:"user"`
}

type VerifyUserEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Token string `json:"verify_token" validate:"required,len=6,numeric"`
}

type VerifyUserEmailResponse struct {
	UserID     string `json:"user_id"`
	IsVerified bool   `json:"is_verified"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ForgotPasswordResponse struct {
	Success bool `json:"success"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Token       string `json:"token" validate:"required,len=6,numeric"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

type ResetPasswordResponse struct {
	Success bool `json:"success"`
}

// Analytics models

type VehiclesInTransitResponse struct {
	Count int64 `json:"count"`
}

type TripsTodayResponse struct {
	Count int64 `json:"count"`
}

// org_models.go

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type OrganizationResponse struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	OwnerUserID    string `json:"owner_user_id"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

type ListUserOrganizationsRequest struct {
	OwnerUserID string `json:"owner_user_id"`
}

type ListUserOrganizationsResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
}

type GetOrganizationMembersListRequest struct {
	OrganizationID string `json:"organization_id"`
}

type GetOrganizationRequest struct {
	OrganizationID string `json:"organization_id"`
}

type OrganizationMemberResponse struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	JoinedAt  int64  `json:"joined_at"`
}

type GetOrganizationMembersListResponse struct {
	Members []OrganizationMemberResponse `json:"members"`
}

type UserIDsRequest struct {
	UserIDs []string
}

type GetUserByEmailRequest struct {
	Email string
}

type AddOrganizationMemberRequest struct {
	OrganizationID string
	UserID         string
	Role           string
}

type InviteUserRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Role  string `json:"role" validate:"required"`
}

type InviteUserResponse struct {
	Success bool `json:"success"`
}

type RemoveOrganizationMemberRequest struct {
	OrganizationID string
	UserID         string
}

type RemoveOrganizationMemberResponse struct {
	Success bool `json:"success"`
}

type UpdateMemberRoleBody struct {
	Role string `json:"role" validate:"required"`
}

type UpdateOrganizationMemberRoleRequest struct {
	OrganizationID string
	UserID         string
	Role           string `json:"role" validate:"required"`
}

type UpdateOrganizationMemberRoleResponse struct {
	Success bool `json:"success"`
}

type CreateVehicleRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
	PlateNumber    string `json:"plate_number" validate:"required,max=50"`
	VehicleType    string `json:"vehicle_type" validate:"max=50"`
	Capacity       int32  `json:"capacity" validate:"gte=0"`
}

type VehicleResponse struct {
	VehicleID      string `json:"vehicle_id"`
	OrganizationID string `json:"organization_id"`
	PlateNumber    string `json:"plate_number"`
	VehicleType    string `json:"vehicle_type" validate:"required,oneof=truck van bike car"`
	Capacity       int32  `json:"capacity"`
	Status         string `json:"status"`
	CreatedAt      int64  `json:"created_at"`
}

type ListVehiclesRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
}

type ListVehiclesResponse struct {
	Vehicles []VehicleResponse `json:"vehicles"`
}

type GetVehicleRequest struct {
	VehicleID string `json:"vehicle_id" validate:"required,uuid"`
}

type CreateDriverRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
	FirstName      string `json:"first_name" validate:"required,max=255"`
	LastName       string `json:"last_name" validate:"required,max=255"`
	VehicleID      string `json:"vehicle_id,omitempty"`
}

type DriverResponse struct {
	DriverID       string `json:"driver_id"`
	OrganizationID string `json:"organization_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	VehicleID      string `json:"vehicle_id"`
	Status         string `json:"status"`
	CreatedAt      int64  `json:"created_at"`
}

type ListDriversRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
}

type ListDriversResponse struct {
	Drivers []DriverResponse `json:"drivers"`
}

type UpdateVehicleRequest struct {
	VehicleID   string `json:"vehicle_id" validate:"required,uuid"`
	PlateNumber string `json:"plate_number" validate:"max=50"`
	VehicleType string `json:"vehicle_type" validate:"omitempty,oneof=truck van bike car"`
	Capacity    int32  `json:"capacity" validate:"gte=0"`
}

type UpdateVehicleStatusRequest struct {
	VehicleID string `json:"vehicle_id" validate:"required,uuid"`
	Status    string `json:"status" validate:"required,oneof=active inactive"`
}

type UpdateDriverStatusRequest struct {
	DriverID string `json:"driver_id" validate:"required,uuid"`
	Status   string `json:"status" validate:"required,oneof=active inactive"`
}

type UpdateDriverRequest struct {
	DriverID  string `json:"driver_id" validate:"required,uuid"`
	FirstName string `json:"first_name" validate:"max=255"`
	LastName  string `json:"last_name" validate:"max=255"`
}

type CreateTripRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
	VehicleID      string `json:"vehicle_id" validate:"required,uuid"`
	DriverID       string `json:"driver_id" validate:"required,uuid"`
}

type TripResponse struct {
	TripID         string `json:"trip_id"`
	OrganizationID string `json:"organization_id"`
	VehicleID      string `json:"vehicle_id"`
	DriverID       string `json:"driver_id"`
	Status         string `json:"status"`
	CreatedAt      int64  `json:"created_at"`
}

type ListTripsRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
}

type ListTripsResponse struct {
	Trips []TripResponse `json:"trips"`
}

type StartTripRequest struct {
	TripID string `json:"trip_id" validate:"required,uuid"`
}

type CompleteTripRequest struct {
	TripID string `json:"trip_id" validate:"required,uuid"`
}

type ListAllVehiclesRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type ListAllVehiclesResponse struct {
	Vehicles []VehicleResponse `json:"vehicles"`
}

type ListAllTripsRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// TotalMembersResponse represents total members count
type TotalMembersResponse struct {
	Count uint64 `json:"count"`
}

// ActiveUsersTodayResponse represents active users today count
type ActiveUsersTodayResponse struct {
	Count uint64 `json:"count"`
}

// Event represents a single analytics event
type Event struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	OrgID     string `json:"org_id"`
	Service   string `json:"service"`
	EventTime string `json:"event_time"`
}

// RecentActivityResponse represents activity feed
type RecentActivityResponse struct {
	Events []Event `json:"events"`
}

type GetTotalMembersRequest struct {
	OwnerUserID string `json:"owner_user_id"`
}