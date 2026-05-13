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

// StudentsService описывает операции модуля студентов.
type StudentsService interface {
	// Create создаёт студента от имени администратора.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateStudent) (*dto.Student, error)

	// List возвращает список студентов.
	List(ctx context.Context, actor dto.Actor) ([]dto.Student, error)

	// GetByID возвращает студента по идентификатору.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Student, error)

	// GetMe возвращает студента, связанного с текущим пользователем.
	GetMe(ctx context.Context, actor dto.Actor) (*dto.Student, error)

	// Update обновляет студента от имени администратора.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateStudent) (*dto.Student, error)

	// Delete удаляет студента от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type studentsService struct {
	students    repositories.StudentRepository
	curriculums repositories.CurriculumRepository
	validator   *validator.Validate
	logger      zerolog.Logger
}

// NewStudentService создаёт сервис студентов.
func NewStudentService(
	students repositories.StudentRepository,
	curriculums repositories.CurriculumRepository,
	logger zerolog.Logger,
) StudentsService {
	return &studentsService{
		students:    students,
		curriculums: curriculums,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
		logger:      logger,
	}
}

func (service *studentsService) Create(ctx context.Context, actor dto.Actor, input dto.CreateStudent) (*dto.Student, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "student", "", "student creation denied")
		return nil, ErrForbidden
	}
	input.LastName = normalizeText(input.LastName)
	input.FirstName = normalizeText(input.FirstName)
	input.MiddleName = normalizeOptionalText(input.MiddleName)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	student, err := service.students.Create(ctx, repositories.StudentCreateData{
		ID:         uuid.NewString(),
		UserID:     input.UserID,
		GroupID:    input.GroupID,
		LastName:   input.LastName,
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		logRepositoryError(service.logger, err, "student", "", "student creation failed")
		return nil, err
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("student_id", student.ID).
		Str("user_id", student.UserID).
		Str("group_id", student.GroupID).
		Msg("student created")

	return toStudentOutput(student), nil
}

func (service *studentsService) List(ctx context.Context, actor dto.Actor) ([]dto.Student, error) {
	if actor.Role != repositories.UserRoleAdministrator && actor.Role != repositories.UserRoleTeacher {
		logAccessDenied(service.logger, actor, "student", "", "students list denied")
		return nil, ErrForbidden
	}

	students, err := service.students.List(ctx)
	if err != nil {
		logRepositoryError(service.logger, err, "student", "", "students list failed")
		return nil, err
	}
	if actor.Role == repositories.UserRoleTeacher {
		availableGroupIDs, err := service.teacherGroupIDs(ctx, actor.ID)
		if err != nil {
			return nil, err
		}
		students = filterStudentsByGroup(students, availableGroupIDs)
		service.logger.Info().
			Str("actor_id", actor.ID).
			Int("groups_count", len(availableGroupIDs)).
			Int("students_count", len(students)).
			Msg("teacher students list filtered")
	}

	output := make([]dto.Student, 0, len(students))
	for _, student := range students {
		output = append(output, *toStudentOutput(&student))
	}

	return output, nil
}

func (service *studentsService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Student, error) {
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	student, err := service.students.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}
	if actor.Role == repositories.UserRoleTeacher {
		if err := service.ensureTeacherCanReadStudent(ctx, actor.ID, student); err != nil {
			return nil, err
		}
	} else if actor.Role != repositories.UserRoleAdministrator && actor.ID != student.UserID {
		logAccessDenied(service.logger, actor, "student", student.ID, "student access denied")
		return nil, ErrForbidden
	}

	return toStudentOutput(student), nil
}

func (service *studentsService) ensureTeacherCanReadStudent(ctx context.Context, teacherUserID string, student *repositories.Student) error {
	groupIDs, err := service.teacherGroupIDs(ctx, teacherUserID)
	if err != nil {
		return err
	}
	if _, ok := groupIDs[student.GroupID]; !ok {
		service.logger.Warn().
			Str("actor_id", teacherUserID).
			Str("actor_role", string(repositories.UserRoleTeacher)).
			Str("student_id", student.ID).
			Str("group_id", student.GroupID).
			Msg("teacher student access denied")
		return ErrForbidden
	}

	return nil
}

func (service *studentsService) teacherGroupIDs(ctx context.Context, teacherUserID string) (map[string]struct{}, error) {
	curriculums, err := service.curriculums.List(ctx)
	if err != nil {
		logRepositoryError(service.logger, err, "curriculum", "", "teacher groups lookup failed")
		return nil, err
	}

	groupIDs := make(map[string]struct{})
	for _, curriculum := range curriculums {
		if curriculum.LeadBy == teacherUserID {
			groupIDs[curriculum.GroupID] = struct{}{}
		}
	}

	return groupIDs, nil
}

func filterStudentsByGroup(students []repositories.Student, groupIDs map[string]struct{}) []repositories.Student {
	filtered := make([]repositories.Student, 0, len(students))
	for _, student := range students {
		if _, ok := groupIDs[student.GroupID]; ok {
			filtered = append(filtered, student)
		}
	}

	return filtered
}

func (service *studentsService) GetMe(ctx context.Context, actor dto.Actor) (*dto.Student, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	student, err := service.students.GetByUserID(ctx, actor.ID)
	if err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "student", actor.ID, "current student lookup failed")
		return nil, mappedErr
	}

	return toStudentOutput(student), nil
}

func (service *studentsService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateStudent) (*dto.Student, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "student", id, "student update denied")
		return nil, ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if input.GroupID == nil && input.LastName == nil && input.FirstName == nil && input.MiddleName == nil {
		return nil, fmt.Errorf("%w: at least one field is required", ErrInvalidInput)
	}
	middleNameSet := input.MiddleName != nil
	if input.LastName != nil {
		lastName := normalizeText(*input.LastName)
		input.LastName = &lastName
	}
	if input.FirstName != nil {
		firstName := normalizeText(*input.FirstName)
		input.FirstName = &firstName
	}
	input.MiddleName = normalizeOptionalText(input.MiddleName)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	student, err := service.students.Update(ctx, id, repositories.StudentUpdateData{
		GroupID:       input.GroupID,
		LastName:      input.LastName,
		FirstName:     input.FirstName,
		MiddleName:    input.MiddleName,
		MiddleNameSet: middleNameSet,
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "student", id, "student update failed")
		return nil, mappedErr
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("student_id", student.ID).
		Str("user_id", student.UserID).
		Str("group_id", student.GroupID).
		Msg("student updated")

	return toStudentOutput(student), nil
}

func (service *studentsService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "student", id, "student deletion denied")
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.students.Delete(ctx, id); err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "student", id, "student deletion failed")
		return mappedErr
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("student_id", id).
		Msg("student deleted")

	return nil
}

func normalizeOptionalText(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := normalizeText(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}

func toStudentOutput(student *repositories.Student) *dto.Student {
	return &dto.Student{
		ID:         student.ID,
		UserID:     student.UserID,
		GroupID:    student.GroupID,
		LastName:   student.LastName,
		FirstName:  student.FirstName,
		MiddleName: student.MiddleName,
		CreatedAt:  student.CreatedAt,
		UpdatedAt:  student.UpdatedAt,
	}
}
