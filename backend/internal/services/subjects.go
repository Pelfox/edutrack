package services

import (
	"context"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// SubjectsService описывает операции модуля предметов.
type SubjectsService interface {
	// Create создаёт предмет от имени администратора.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateSubject) (*dto.Subject, error)

	// List возвращает список предметов.
	List(ctx context.Context, actor dto.Actor) ([]dto.Subject, error)

	// GetByID возвращает предмет по идентификатору.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Subject, error)

	// Update обновляет предмет от имени администратора.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateSubject) (*dto.Subject, error)

	// Delete удаляет предмет от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type subjectsService struct {
	subjects  repositories.SubjectRepository
	validator *validator.Validate
	logger    zerolog.Logger
}

// NewSubjectService создаёт сервис предметов.
func NewSubjectService(subjects repositories.SubjectRepository, logger zerolog.Logger) SubjectsService {
	return &subjectsService{
		subjects:  subjects,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		logger:    logger,
	}
}

func (service *subjectsService) Create(ctx context.Context, actor dto.Actor, input dto.CreateSubject) (*dto.Subject, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "subject", "", "subject creation denied")
		return nil, ErrForbidden
	}
	input.Title = normalizeText(input.Title)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	subject, err := service.subjects.Create(ctx, repositories.SubjectCreateData{
		ID:        uuid.NewString(),
		Title:     input.Title,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		logRepositoryError(service.logger, err, "subject", "", "subject creation failed")
		return nil, err
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("subject_id", subject.ID).
		Str("title", subject.Title).
		Msg("subject created")

	return toSubjectOutput(subject), nil
}

func (service *subjectsService) List(ctx context.Context, actor dto.Actor) ([]dto.Subject, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	subjects, err := service.subjects.List(ctx)
	if err != nil {
		logRepositoryError(service.logger, err, "subject", "", "subjects list failed")
		return nil, err
	}

	output := make([]dto.Subject, 0, len(subjects))
	for _, subject := range subjects {
		output = append(output, *toSubjectOutput(&subject))
	}

	return output, nil
}

func (service *subjectsService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.Subject, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	subject, err := service.subjects.GetByID(ctx, id)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toSubjectOutput(subject), nil
}

func (service *subjectsService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateSubject) (*dto.Subject, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "subject", id, "subject update denied")
		return nil, ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	input.Title = normalizeText(input.Title)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	subject, err := service.subjects.Update(ctx, id, repositories.SubjectUpdateData{Title: input.Title})
	if err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "subject", id, "subject update failed")
		return nil, mappedErr
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("subject_id", subject.ID).
		Str("title", subject.Title).
		Msg("subject updated")

	return toSubjectOutput(subject), nil
}

func (service *subjectsService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		logAccessDenied(service.logger, actor, "subject", id, "subject deletion denied")
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.subjects.Delete(ctx, id); err != nil {
		mappedErr := mapDirectoryRepositoryError(err)
		logRepositoryError(service.logger, mappedErr, "subject", id, "subject deletion failed")
		return mappedErr
	}

	service.logger.Info().
		Str("actor_id", actor.ID).
		Str("subject_id", id).
		Msg("subject deleted")

	return nil
}

func toSubjectOutput(subject *repositories.Subject) *dto.Subject {
	return &dto.Subject{
		ID:        subject.ID,
		Title:     subject.Title,
		CreatedAt: subject.CreatedAt,
	}
}
