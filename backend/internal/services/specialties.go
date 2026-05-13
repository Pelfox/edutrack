package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// SpecialtiesService описывает операции модуля специальностей.
type SpecialtiesService interface {
	// Create создаёт специальность от имени администратора.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateSpecialty) (*dto.Specialty, error)

	// List возвращает список специальностей.
	List(ctx context.Context, actor dto.Actor) ([]dto.Specialty, error)

	// GetByID возвращает специальность по идентификатору.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Specialty, error)

	// Update обновляет специальность от имени администратора.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateSpecialty) (*dto.Specialty, error)

	// Delete удаляет специальность от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type specialtiesService struct {
	specialties repositories.SpecialtyRepository
	validator   *validator.Validate
	logger      zerolog.Logger
}

// NewSpecialtyService создаёт сервис специальностей.
func NewSpecialtyService(specialties repositories.SpecialtyRepository, logger zerolog.Logger) SpecialtiesService {
	return &specialtiesService{
		specialties: specialties,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
		logger:      logger,
	}
}

func (service *specialtiesService) Create(ctx context.Context, actor dto.Actor, input dto.CreateSpecialty) (*dto.Specialty, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "specialty", "", "specialty creation denied")
		return nil, ErrForbidden
	}
	input.Title = normalizeText(input.Title)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	specialty, err := service.specialties.Create(ctx, repositories.SpecialtyCreateData{
		ID:        uuid.NewString(),
		Title:     input.Title,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		logRepositoryError(service.logger, err, "specialty", "", "specialty creation failed")
		return nil, err
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("specialty_id", specialty.ID).
		Str("title", specialty.Title).
		Msg("specialty created")

	return toSpecialtyOutput(specialty), nil
}

func (service *specialtiesService) List(ctx context.Context, actor dto.Actor) ([]dto.Specialty, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	specialties, err := service.specialties.List(ctx)
	if err != nil {
		logRepositoryError(service.logger, err, "specialty", "", "specialties list failed")
		return nil, err
	}

	output := make([]dto.Specialty, 0, len(specialties))
	for _, specialty := range specialties {
		output = append(output, *toSpecialtyOutput(&specialty))
	}

	return output, nil
}

func (service *specialtiesService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Specialty, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	specialty, err := service.specialties.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toSpecialtyOutput(specialty), nil
}

func (service *specialtiesService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateSpecialty) (*dto.Specialty, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "specialty", id, "specialty update denied")
		return nil, ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	input.Title = normalizeText(input.Title)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	specialty, err := service.specialties.Update(ctx, id, repositories.SpecialtyUpdateData{
		Title: input.Title,
	})
	if err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "specialty", id, "specialty update failed")
		return nil, mappedErr
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("specialty_id", specialty.ID).
		Str("title", specialty.Title).
		Msg("specialty updated")

	return toSpecialtyOutput(specialty), nil
}

func (service *specialtiesService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "specialty", id, "specialty deletion denied")
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.specialties.Delete(ctx, id); err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "specialty", id, "specialty deletion failed")
		return mappedErr
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("specialty_id", id).
		Msg("specialty deleted")

	return nil
}

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func mapDirectoryRepositoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func toSpecialtyOutput(specialty *repositories.Specialty) *dto.Specialty {
	return &dto.Specialty{
		ID:        specialty.ID,
		Title:     specialty.Title,
		CreatedAt: specialty.CreatedAt,
	}
}
