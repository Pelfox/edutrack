package dto

import "time"

// Grade содержит представление оценки для HTTP-ответов.
type Grade struct {
	ID           string    `json:"id"`
	CurriculumID string    `json:"curriculum_id"`
	StudentID    string    `json:"student_id"`
	AuthorID     string    `json:"author_id"`
	Value        int       `json:"value"`
	Comment      *string   `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateGrade содержит данные для создания оценки.
type CreateGrade struct {
	CurriculumID string  `json:"curriculum_id" binding:"required" validate:"required,uuid4"`
	StudentID    string  `json:"student_id" binding:"required" validate:"required,uuid4"`
	Value        *int    `json:"value" binding:"required" validate:"required"`
	Comment      *string `json:"comment" validate:"omitempty,max=255"`
}

// UpdateGrade содержит данные для частичного обновления оценки.
type UpdateGrade struct {
	CurriculumID *string `json:"curriculum_id" validate:"omitempty,uuid4"`
	StudentID    *string `json:"student_id" validate:"omitempty,uuid4"`
	Value        *int    `json:"value"`
	Comment      *string `json:"comment" validate:"omitempty,max=255"`
}
