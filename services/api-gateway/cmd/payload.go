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