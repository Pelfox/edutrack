package services

import (
	"context"
	"fmt"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CurriculumsService описывает операции модуля учебных планов.
type CurriculumsService interface {
	// Create создаёт учебный план от имени администратора.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateCurriculum) (*dto.Curriculum, error)

	// List возвращает список учебных планов.
	List(ctx context.Context, actor dto.Actor) ([]dto.Curriculum, error)

	// GetByID возвращает учебный план по идентификатору.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Curriculum, error)

	// Update обновляет учебный план от имени администратора.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateCurriculum) (*dto.Curriculum, error)

	// Delete удаляет учебный план от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type curriculumsService struct {
	curriculums repositories.CurriculumRepository
	validator   *validator.Validate
}

// NewCurriculumService создаёт сервис учебных планов.
func NewCurriculumService(curriculums repositories.CurriculumRepository) CurriculumsService {
	return &curriculumsService{
		curriculums: curriculums,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (service *curriculumsService) Create(ctx context.Context, actor dto.Actor, input dto.CreateCurriculum) (*dto.Curriculum, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	curriculum, err := service.curriculums.Create(ctx, repositories.CurriculumCreateData{
		ID:         uuid.NewString(),
		Hours:      input.Hours,
		Semester:   input.Semester,
		ReportType: input.ReportType,
		SubjectID:  input.SubjectID,
		GroupID:    input.GroupID,
		LeadBy:     input.LeadBy,
	})
	if err != nil {
		return nil, err
	}

	return toCurriculumOutput(curriculum), nil
}

func (service *curriculumsService) List(ctx context.Context, actor dto.Actor) ([]dto.Curriculum, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	curriculums, err := service.curriculums.List(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Curriculum, 0, len(curriculums))
	for _, curriculum := range curriculums {
		output = append(output, *toCurriculumOutput(&curriculum))
	}

	return output, nil
}

func (service *curriculumsService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Curriculum, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	curriculum, err := service.curriculums.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toCurriculumOutput(curriculum), nil
}

func (service *curriculumsService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateCurriculum) (*dto.Curriculum, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if input.Hours == nil && input.Semester == nil && input.ReportType == nil && input.SubjectID == nil && input.GroupID == nil && input.LeadBy == nil {
		return nil, fmt.Errorf("%w: at least one field is required", ErrInvalidInput)
	}
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	curriculum, err := service.curriculums.Update(ctx, id, repositories.CurriculumUpdateData{
		Hours:      input.Hours,
		Semester:   input.Semester,
		ReportType: input.ReportType,
		SubjectID:  input.SubjectID,
		GroupID:    input.GroupID,
		LeadBy:     input.LeadBy,
	})
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toCurriculumOutput(curriculum), nil
}

func (service *curriculumsService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.curriculums.Delete(ctx, id); err != nil {
		return mapDirectoryRepositoryError(err)
	}

	return nil
}

func toCurriculumOutput(curriculum *repositories.Curriculum) *dto.Curriculum {
	return &dto.Curriculum{
		ID:         curriculum.ID,
		Hours:      curriculum.Hours,
		Semester:   curriculum.Semester,
		ReportType: curriculum.ReportType,
		SubjectID:  curriculum.SubjectID,
		GroupID:    curriculum.GroupID,
		LeadBy:     curriculum.LeadBy,
	}
}
