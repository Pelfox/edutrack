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
	students  repositories.StudentRepository
	validator *validator.Validate
	now       func() time.Time
}

// NewStudentService создаёт сервис студентов.
func NewStudentService(students repositories.StudentRepository) StudentsService {
	return &studentsService{
		students:  students,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		now:       time.Now,
	}
}

func (service *studentsService) Create(ctx context.Context, actor dto.Actor, input dto.CreateStudent) (*dto.Student, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	input.LastName = normalizeText(input.LastName)
	input.FirstName = normalizeText(input.FirstName)
	input.MiddleName = normalizeOptionalText(input.MiddleName)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	now := service.now().UTC()
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
		return nil, err
	}

	return toStudentOutput(student), nil
}

func (service *studentsService) List(ctx context.Context, actor dto.Actor) ([]dto.Student, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}

	students, err := service.students.List(ctx)
	if err != nil {
		return nil, err
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
	if actor.Role != repositories.UserRoleAdministrator && actor.ID != student.UserID {
		return nil, ErrForbidden
	}

	return toStudentOutput(student), nil
}

func (service *studentsService) GetMe(ctx context.Context, actor dto.Actor) (*dto.Student, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	student, err := service.students.GetByUserID(ctx, actor.ID)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toStudentOutput(student), nil
}

func (service *studentsService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateStudent) (*dto.Student, error) {
	if actor.Role != repositories.UserRoleAdministrator {
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
		UpdatedAt:     service.now().UTC(),
	})
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toStudentOutput(student), nil
}

func (service *studentsService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.students.Delete(ctx, id); err != nil {
		return mapDirectoryRepositoryError(err)
	}

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
