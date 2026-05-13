package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// GradesService описывает операции модуля оценок.
type GradesService interface {
	// Create создаёт оценку от имени администратора или преподавателя.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateGrade) (*dto.Grade, error)

	// List возвращает список оценок с учётом роли пользователя.
	List(ctx context.Context, actor dto.Actor) ([]dto.Grade, error)

	// GetByID возвращает оценку по идентификатору с учётом роли пользователя.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Grade, error)

	// Update обновляет оценку от имени администратора или преподавателя.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateGrade) (*dto.Grade, error)

	// Delete удаляет оценку от имени администратора или преподавателя.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type gradesService struct {
	grades      repositories.GradeRepository
	students    repositories.StudentRepository
	curriculums repositories.CurriculumRepository
	validator   *validator.Validate
	logger      zerolog.Logger
}

// NewGradeService создаёт сервис оценок.
func NewGradeService(
	grades repositories.GradeRepository,
	students repositories.StudentRepository,
	curriculums repositories.CurriculumRepository,
	logger zerolog.Logger,
) GradesService {
	return &gradesService{
		grades:      grades,
		students:    students,
		curriculums: curriculums,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
		logger:      logger,
	}
}

func (service *gradesService) Create(ctx context.Context, actor dto.Actor, input dto.CreateGrade) (*dto.Grade, error) {
	if !canManageGrades(actor) {
		service.logGradeAccessDenied(actor, "", "grade creation denied")
		return nil, ErrForbidden
	}
	input.Comment = normalizeOptionalText(input.Comment)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}
	if err := service.ensureCanWriteGrade(ctx, actor, input.CurriculumID, input.StudentID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	grade, err := service.grades.Create(ctx, repositories.GradeCreateData{
		ID:           uuid.NewString(),
		CurriculumID: input.CurriculumID,
		StudentID:    input.StudentID,
		AuthorID:     actor.ID,
		Value:        *input.Value,
		Comment:      input.Comment,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("actor_role", string(actor.Role)).
		Str("grade_id", grade.ID).
		Str("curriculum_id", grade.CurriculumID).
		Str("student_id", grade.StudentID).
		Int("value", grade.Value).
		Msg("grade created")

	return toGradeOutput(grade), nil
}

func (service *gradesService) List(ctx context.Context, actor dto.Actor) ([]dto.Grade, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	var (
		grades []repositories.Grade
		err    error
	)
	if actor.Role == repositories.UserRoleStudent {
		student, err := service.students.GetByUserID(ctx, actor.ID)
		if err != nil {
			return nil, mapDirectoryRepositoryError(err)
		}
		grades, err = service.grades.ListByStudentID(ctx, student.ID)
	} else {
		grades, err = service.grades.List(ctx)
	}
	if err != nil {
		return nil, err
	}

	output := make([]dto.Grade, 0, len(grades))
	for _, grade := range grades {
		output = append(output, *toGradeOutput(&grade))
	}

	return output, nil
}

func (service *gradesService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Grade, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	grade, err := service.grades.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}
	if err := service.ensureCanReadGrade(ctx, actor, grade); err != nil {
		return nil, err
	}

	return toGradeOutput(grade), nil
}

func (service *gradesService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateGrade) (*dto.Grade, error) {
	if !canManageGrades(actor) {
		service.logGradeAccessDenied(actor, id, "grade update denied")
		return nil, ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if input.CurriculumID == nil && input.StudentID == nil && input.Value == nil && input.Comment == nil {
		return nil, fmt.Errorf("%w: at least one field is required", ErrInvalidInput)
	}
	commentSet := input.Comment != nil
	input.Comment = normalizeOptionalText(input.Comment)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}
	existingGrade, err := service.grades.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}
	curriculumID := existingGrade.CurriculumID
	if input.CurriculumID != nil {
		curriculumID = *input.CurriculumID
	}
	studentID := existingGrade.StudentID
	if input.StudentID != nil {
		studentID = *input.StudentID
	}
	if err := service.ensureCanWriteGrade(ctx, actor, curriculumID, studentID); err != nil {
		return nil, err
	}

	grade, err := service.grades.Update(ctx, id, repositories.GradeUpdateData{
		CurriculumID: input.CurriculumID,
		StudentID:    input.StudentID,
		Value:        input.Value,
		Comment:      input.Comment,
		CommentSet:   commentSet,
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("actor_role", string(actor.Role)).
		Str("grade_id", grade.ID).
		Str("curriculum_id", grade.CurriculumID).
		Str("student_id", grade.StudentID).
		Int("value", grade.Value).
		Msg("grade updated")

	return toGradeOutput(grade), nil
}

func (service *gradesService) ensureCanWriteGrade(ctx context.Context, actor dto.Actor, curriculumID string, studentID string) error {
	if actor.Role == repositories.UserRoleAdministrator {
		return nil
	}

	curriculum, err := service.curriculums.GetByID(ctx, curriculumID)
	if err != nil {
		return mapDirectoryRepositoryError(err)
	}
	if curriculum.LeadBy != actor.ID {
		service.logger.Warn().
			Str("actor_id", actor.ID).
			Str("actor_role", string(actor.Role)).
			Str("curriculum_id", curriculumID).
			Str("student_id", studentID).
			Str("lead_by", curriculum.LeadBy).
			Msg("grade write denied")
		return ErrForbidden
	}

	student, err := service.students.GetByID(ctx, studentID)
	if err != nil {
		return mapDirectoryRepositoryError(err)
	}
	if student.GroupID != curriculum.GroupID {
		return fmt.Errorf("%w: student is not assigned to curriculum group", ErrInvalidInput)
	}

	return nil
}

func (service *gradesService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if !canManageGrades(actor) {
		service.logGradeAccessDenied(actor, id, "grade deletion denied")
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	grade, err := service.grades.GetByID(ctx, id)
	if err != nil {
		return mapDirectoryRepositoryError(err)
	}
	if err := service.ensureCanWriteGrade(ctx, actor, grade.CurriculumID, grade.StudentID); err != nil {
		return err
	}

	if err := service.grades.Delete(ctx, id); err != nil {
		return mapDirectoryRepositoryError(err)
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("actor_role", string(actor.Role)).
		Str("grade_id", id).
		Msg("grade deleted")

	return nil
}

func (service *gradesService) ensureCanReadGrade(ctx context.Context, actor dto.Actor, grade *repositories.Grade) error {
	if actor.Role != repositories.UserRoleStudent {
		return nil
	}

	student, err := service.students.GetByUserID(ctx, actor.ID)
	if err != nil {
		return mapDirectoryRepositoryError(err)
	}
	if grade.StudentID != student.ID {
		service.logger.Warn().
			Str("actor_id", actor.ID).
			Str("actor_role", string(actor.Role)).
			Str("grade_id", grade.ID).
			Str("student_id", student.ID).
			Msg("grade read denied")
		return ErrForbidden
	}

	return nil
}

func (service *gradesService) logGradeAccessDenied(actor dto.Actor, gradeID string, message string) {
	event := service.logger.Warn().
		Str("actor_id", actor.ID).
		Str("actor_role", string(actor.Role))
	if gradeID != "" {
		event = event.Str("grade_id", gradeID)
	}

	event.Msg(message)
}

func canManageGrades(actor dto.Actor) bool {
	return actor.Role == repositories.UserRoleAdministrator || actor.Role == repositories.UserRoleTeacher
}

func toGradeOutput(grade *repositories.Grade) *dto.Grade {
	return &dto.Grade{
		ID:           grade.ID,
		CurriculumID: grade.CurriculumID,
		StudentID:    grade.StudentID,
		AuthorID:     grade.AuthorID,
		Value:        grade.Value,
		Comment:      grade.Comment,
		CreatedAt:    grade.CreatedAt,
		UpdatedAt:    grade.UpdatedAt,
	}
}
