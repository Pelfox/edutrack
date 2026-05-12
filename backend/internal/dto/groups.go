package dto

import (
	"time"

	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

// Group содержит представление группы для HTTP-ответов.
type Group struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	StudyForm     repositories.StudyForm `json:"study_form"`
	AdmissionYear int                    `json:"admission_year"`
	SpecialtyID   string                 `json:"specialty_id"`
	CreatedAt     time.Time              `json:"created_at"`
}

// CreateGroup содержит данные для создания группы.
type CreateGroup struct {
	Name          string                 `json:"name" binding:"required" validate:"required,min=1,max=255"`
	StudyForm     repositories.StudyForm `json:"study_form" binding:"required" validate:"required,oneof=full_time evening extramural"`
	AdmissionYear int                    `json:"admission_year" binding:"required" validate:"required,min=1900,max=2100"`
	SpecialtyID   string                 `json:"specialty_id" binding:"required" validate:"required,uuid4"`
}

// UpdateGroup содержит данные для частичного обновления группы.
type UpdateGroup struct {
	Name          *string                 `json:"name" validate:"omitempty,min=1,max=255"`
	StudyForm     *repositories.StudyForm `json:"study_form" validate:"omitempty,oneof=full_time evening extramural"`
	AdmissionYear *int                    `json:"admission_year" validate:"omitempty,min=1900,max=2100"`
	SpecialtyID   *string                 `json:"specialty_id" validate:"omitempty,uuid4"`
}
