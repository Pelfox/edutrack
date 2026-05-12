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
