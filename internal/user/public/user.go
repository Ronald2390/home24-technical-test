package public

import "home24-technical-test/internal/user/model"

// CreateUserParams represents object to create user
type CreateUserParams struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Password string `json:"password"`
}

//FindAllUsersParams params for find all
type FindAllUsersParams struct {
	Page   int    `json:"page"`
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// LoginResponse represents the response of login function
type LoginResponse struct {
	SessionID string      `json:"sessionId"`
	User      *model.User `json:"user"`
}

// LoginParams represent the http request data for login user
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChangePasswordParams represent the http request data for change password
// swagger:model
type ChangePasswordParams struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// UpdateUserParams represent the http request data for update user
type UpdateUserParams struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Password string `json:"password"`
}
