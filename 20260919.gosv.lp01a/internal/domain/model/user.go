package model

import "time"

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=3"`
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=3,max=20,startswith=@"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserUpdateRequest struct {
	UserID          string  `json:"-"`
	Name            *string `json:"name" validate:"omitempty,min=3"`
	Email           *string `json:"email" validate:"omitempty,email"`
	Username        *string `json:"username" validate:"omitempty,min=3,max=20,startswith=@"`
	Password        *string `json:"password" validate:"omitempty,min=8"`
	ConfirmPassword *string `json:"confirmPassword" validate:"omitempty,eqfield=Password"`
	OldPassword     *string `json:"oldPassword" validate:"omitempty"`
}

type UserDeleteRequest struct {
	UserID       string `json:"-"`
	Password     string `json:"password" validate:"required"`
	Confirmation string `json:"confirmation" validate:"required,oneof=DELETE"`
}

type UserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
