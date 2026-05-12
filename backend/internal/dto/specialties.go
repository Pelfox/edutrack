package dto

import "time"

// Specialty содержит представление специальности для HTTP-ответов.
type Specialty struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateSpecialty содержит данные для создания специальности.
type CreateSpecialty struct {
	Title string `json:"title" binding:"required" validate:"required,min=1,max=255"`
}

// UpdateSpecialty содержит данные для обновления специальности.
type UpdateSpecialty struct {
	Title string `json:"title" binding:"required" validate:"required,min=1,max=255"`
}
