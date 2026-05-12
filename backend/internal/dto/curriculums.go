package dto

import "github.com/Pelfox/edutrack/backend/internal/repositories"

// Curriculum содержит представление учебного плана для HTTP-ответов.
type Curriculum struct {
	ID         string                            `json:"id"`
	Hours      int                               `json:"hours"`
	Semester   int                               `json:"semester"`
	ReportType repositories.CurriculumReportType `json:"report_type"`
	SubjectID  string                            `json:"subject_id"`
	GroupID    string                            `json:"group_id"`
	LeadBy     string                            `json:"lead_by"`
}

// CreateCurriculum содержит данные для создания учебного плана.
type CreateCurriculum struct {
	Hours      int                               `json:"hours" binding:"required" validate:"required,min=1"`
	Semester   int                               `json:"semester" binding:"required" validate:"required,min=1,max=12"`
	ReportType repositories.CurriculumReportType `json:"report_type" binding:"required" validate:"required,oneof=exam test diff_test"`
	SubjectID  string                            `json:"subject_id" binding:"required" validate:"required,uuid4"`
	GroupID    string                            `json:"group_id" binding:"required" validate:"required,uuid4"`
	LeadBy     string                            `json:"lead_by" binding:"required" validate:"required,uuid4"`
}

// UpdateCurriculum содержит данные для частичного обновления учебного плана.
type UpdateCurriculum struct {
	Hours      *int                               `json:"hours" validate:"omitempty,min=1"`
	Semester   *int                               `json:"semester" validate:"omitempty,min=1,max=12"`
	ReportType *repositories.CurriculumReportType `json:"report_type" validate:"omitempty,oneof=exam test diff_test"`
	SubjectID  *string                            `json:"subject_id" validate:"omitempty,uuid4"`
	GroupID    *string                            `json:"group_id" validate:"omitempty,uuid4"`
	LeadBy     *string                            `json:"lead_by" validate:"omitempty,uuid4"`
}
