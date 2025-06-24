package user

import (
	"github.com/mukezhz/appointment-booking/domain/constants"
	"github.com/mukezhz/appointment-booking/domain/models"
	"github.com/mukezhz/appointment-booking/pkg/types"
)

// RegisterUserDTO represents the data required for user registration
type RegisterUserDTO struct {
	Email      string `json:"email" binding:"required,email"`
	Role       string `json:"role" binding:"required,oneof=admin doctor patient"`
	FullName   string `json:"full_name" binding:"required"`
	FullNameJa string `json:"full_name_ja,omitempty"`
}

// LoginUserDTO represents the data required for user login
type LoginUserDTO struct {
	Email string `json:"email" binding:"required,email"`
}

// UserResponseDTO represents the user data returned to clients
type UserResponseDTO struct {
	UUID       types.BinaryUUID   `json:"uuid"`
	Email      string             `json:"email"`
	Role       constants.UserRole `json:"role"`
	FullName   string             `json:"full_name"`
	FullNameJa string             `json:"full_name_ja,omitempty"`
	IsActive   bool               `json:"is_active"`
}

// UpdateUserDTO represents the data that can be updated for a user
type UpdateUserDTO struct {
	FullName   string `json:"full_name,omitempty"`
	FullNameJa string `json:"full_name_ja,omitempty"`
}

func toUserResponseDTO(user *models.User) UserResponseDTO {
	return UserResponseDTO{
		UUID:       user.UUID,
		Email:      user.Email,
		Role:       user.Role,
		FullName:   user.FullName,
		FullNameJa: user.FullNameJa,
		IsActive:   user.IsActive,
	}
}
