package dto

import (
	"time"

	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

// Actor содержит данные авторизованного пользователя, необходимые для проверки доступа.
type Actor struct {
	ID   string                `validate:"required,uuid4"`
	Role repositories.UserRole `validate:"required,oneof=administrator teacher student"`
}

// User содержит безопасное представление пользователя для HTTP-ответов.
type User struct {
	ID        string                `json:"id"`
	Email     string                `json:"email"`
	Role      repositories.UserRole `json:"role"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// CreateUser содержит данные для создания нового пользователя.
type CreateUser struct {
	Email    string                `json:"email" binding:"required" validate:"required,email"`
	Password string                `json:"password" binding:"required" validate:"required,min=8"`
	Role     repositories.UserRole `json:"role" binding:"required" validate:"required,oneof=administrator teacher student"`
}

// UpdateUser содержит данные для частичного обновления пользователя.
type UpdateUser struct {
	Email    *string                `json:"email" validate:"omitempty,email"`
	Password *string                `json:"password" validate:"omitempty,min=8"`
	Role     *repositories.UserRole `json:"role" validate:"omitempty,oneof=administrator teacher student"`
}
