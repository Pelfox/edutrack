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

// GroupsService описывает операции модуля групп.
type GroupsService interface {
	// Create создаёт группу от имени администратора.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateGroup) (*dto.Group, error)

	// List возвращает список групп.
	List(ctx context.Context, actor dto.Actor) ([]dto.Group, error)

	// GetByID возвращает группу по идентификатору.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Group, error)

	// Update обновляет группу от имени администратора.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateGroup) (*dto.Group, error)

	// Delete удаляет группу от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type groupsService struct {
	groups    repositories.GroupRepository
	validator *validator.Validate
}

// NewGroupService создаёт сервис групп.
func NewGroupService(groups repositories.GroupRepository) GroupsService {
	return &groupsService{
		groups:    groups,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (service *groupsService) Create(ctx context.Context, actor dto.Actor, input dto.CreateGroup) (*dto.Group, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	input.Name = normalizeText(input.Name)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	group, err := service.groups.Create(ctx, repositories.GroupCreateData{
		ID:            uuid.NewString(),
		Name:          input.Name,
		StudyForm:     input.StudyForm,
		AdmissionYear: input.AdmissionYear,
		SpecialtyID:   input.SpecialtyID,
		CreatedAt:     time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	return toGroupOutput(group), nil
}

func (service *groupsService) List(ctx context.Context, actor dto.Actor) ([]dto.Group, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	groups, err := service.groups.List(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Group, 0, len(groups))
	for _, group := range groups {
		output = append(output, *toGroupOutput(&group))
	}

	return output, nil
}

func (service *groupsService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Group, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	group, err := service.groups.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toGroupOutput(group), nil
}

func (service *groupsService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateGroup) (*dto.Group, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if input.Name == nil && input.StudyForm == nil && input.AdmissionYear == nil && input.SpecialtyID == nil {
		return nil, fmt.Errorf("%w: at least one field is required", ErrInvalidInput)
	}
	if input.Name != nil {
		name := normalizeText(*input.Name)
		input.Name = &name
	}
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	group, err := service.groups.Update(ctx, id, repositories.GroupUpdateData{
		Name:          input.Name,
		StudyForm:     input.StudyForm,
		AdmissionYear: input.AdmissionYear,
		SpecialtyID:   input.SpecialtyID,
	})
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toGroupOutput(group), nil
}

func (service *groupsService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.groups.Delete(ctx, id); err != nil {
		return mapDirectoryRepositoryError(err)
	}

	return nil
}

func toGroupOutput(group *repositories.Group) *dto.Group {
	return &dto.Group{
		ID:            group.ID,
		Name:          group.Name,
		StudyForm:     group.StudyForm,
		AdmissionYear: group.AdmissionYear,
		SpecialtyID:   group.SpecialtyID,
		CreatedAt:     group.CreatedAt,
	}
}
