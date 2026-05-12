package dto

import "time"

// Subject содержит представление предмета для HTTP-ответов.
type Subject struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateSubject содержит данные для создания предмета.
type CreateSubject struct {
	Title string `json:"title" binding:"required" validate:"required,min=1,max=255"`
}

// UpdateSubject содержит данные для обновления предмета.
type UpdateSubject struct {
	Title string `json:"title" binding:"required" validate:"required,min=1,max=255"`
}
