package main

type CreateUserRequest struct {
	FirstName string `json:"firstname" validate:"required,max=255"`
	LastName  string `json:"lastname" validate:"required,max=255"`
	Email     string `json:"email" validate:"max=255"`
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

type LoginResponse struct {
	
}
