package dto

import "time"

// Student содержит представление студента для HTTP-ответов.
type Student struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	GroupID    string    `json:"group_id"`
	LastName   string    `json:"last_name"`
	FirstName  string    `json:"first_name"`
	MiddleName *string   `json:"middle_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateStudent содержит данные для создания студента.
type CreateStudent struct {
	UserID     string  `json:"user_id" binding:"required" validate:"required,uuid4"`
	GroupID    string  `json:"group_id" binding:"required" validate:"required,uuid4"`
	LastName   string  `json:"last_name" binding:"required" validate:"required,min=1,max=255"`
	FirstName  string  `json:"first_name" binding:"required" validate:"required,min=1,max=255"`
	MiddleName *string `json:"middle_name" validate:"omitempty,min=1,max=255"`
}

// UpdateStudent содержит данные для частичного обновления студента.
type UpdateStudent struct {
	GroupID    *string `json:"group_id" validate:"omitempty,uuid4"`
	LastName   *string `json:"last_name" validate:"omitempty,min=1,max=255"`
	FirstName  *string `json:"first_name" validate:"omitempty,min=1,max=255"`
	MiddleName *string `json:"middle_name" validate:"omitempty,min=1,max=255"`
}
