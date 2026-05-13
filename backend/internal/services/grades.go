package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
}

// NewGradeService создаёт сервис оценок.
func NewGradeService(
	grades repositories.GradeRepository,
	students repositories.StudentRepository,
	curriculums repositories.CurriculumRepository,
) GradesService {
	return &gradesService{
		grades:      grades,
		students:    students,
		curriculums: curriculums,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (service *gradesService) Create(ctx context.Context, actor dto.Actor, input dto.CreateGrade) (*dto.Grade, error) {
	if !canManageGrades(actor) {
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
		return ErrForbidden
	}

	return nil
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
