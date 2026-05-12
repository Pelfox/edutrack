package dto

import "github.com/Pelfox/edutrack/backend/internal/repositories"

// Profile содержит ФИО и учётные данные текущего пользователя.
type Profile struct {
	ID         string                `json:"id"`
	UserID     string                `json:"user_id"`
	Email      string                `json:"email"`
	Role       repositories.UserRole `json:"role"`
	LastName   string                `json:"last_name"`
	FirstName  string                `json:"first_name"`
	MiddleName *string               `json:"middle_name"`
}

// CreateProfile содержит данные для создания профильной записи пользователя.
type CreateProfile struct {
	UserID     string  `json:"user_id" binding:"required" validate:"required,uuid4"`
	LastName   string  `json:"last_name" binding:"required" validate:"required,min=1,max=255"`
	FirstName  string  `json:"first_name" binding:"required" validate:"required,min=1,max=255"`
	MiddleName *string `json:"middle_name" validate:"omitempty,min=1,max=255"`
}

// UpdateProfile содержит данные для частичного обновления профильной записи пользователя.
type UpdateProfile struct {
	LastName   *string `json:"last_name" validate:"omitempty,min=1,max=255"`
	FirstName  *string `json:"first_name" validate:"omitempty,min=1,max=255"`
	MiddleName *string `json:"middle_name" validate:"omitempty,min=1,max=255"`
}
