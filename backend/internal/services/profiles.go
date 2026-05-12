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

// ProfileService описывает операции с профилем текущего пользователя.
type ProfileService interface {
	// GetMe возвращает профиль текущего пользователя.
	GetMe(ctx context.Context, actor dto.Actor) (*dto.Profile, error)
}

// StaffProfilesService описывает операции с профилями сотрудников.
type StaffProfilesService interface {
	// Create создаёт профиль сотрудника от имени администратора.
	Create(ctx context.Context, actor dto.Actor, role repositories.UserRole, input dto.CreateProfile) (*dto.Profile, error)

	// List возвращает список профилей сотрудников с указанной ролью.
	List(ctx context.Context, actor dto.Actor, role repositories.UserRole) ([]dto.Profile, error)

	// GetByID возвращает профиль сотрудника по идентификатору.
	GetByID(ctx context.Context, actor dto.Actor, role repositories.UserRole, id string) (*dto.Profile, error)

	// Update обновляет профиль сотрудника от имени администратора.
	Update(ctx context.Context, actor dto.Actor, role repositories.UserRole, id string, input dto.UpdateProfile) (*dto.Profile, error)

	// Delete удаляет профиль сотрудника от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, role repositories.UserRole, id string) error
}

type profileService struct {
	profiles repositories.ProfileRepository
}

type staffProfilesService struct {
	profiles  repositories.ProfileRepository
	users     repositories.UserRepository
	validator *validator.Validate
}

// NewProfileService создаёт сервис профилей пользователей.
func NewProfileService(profiles repositories.ProfileRepository) ProfileService {
	return &profileService{profiles: profiles}
}

// NewStaffProfilesService создаёт сервис профилей сотрудников.
func NewStaffProfilesService(profiles repositories.ProfileRepository, users repositories.UserRepository) StaffProfilesService {
	return &staffProfilesService{
		profiles:  profiles,
		users:     users,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (service *profileService) GetMe(ctx context.Context, actor dto.Actor) (*dto.Profile, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	profile, err := service.profiles.GetByUserID(ctx, actor.ID, actor.Role)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toProfileOutput(profile), nil
}

func (service *staffProfilesService) Create(ctx context.Context, actor dto.Actor, role repositories.UserRole, input dto.CreateProfile) (*dto.Profile, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateStaffProfileRole(role); err != nil {
		return nil, err
	}
	input.LastName = normalizeText(input.LastName)
	input.FirstName = normalizeText(input.FirstName)
	input.MiddleName = normalizeOptionalText(input.MiddleName)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	user, err := service.users.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, mapUserRepositoryError(err)
	}
	if user.Role != role {
		return nil, fmt.Errorf("%w: user role does not match profile role", ErrInvalidInput)
	}

	now := time.Now().UTC()
	profile, err := service.profiles.Create(ctx, role, repositories.ProfileCreateData{
		ID:         uuid.NewString(),
		UserID:     input.UserID,
		LastName:   input.LastName,
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return nil, err
	}

	return toProfileOutput(profile), nil
}

func (service *staffProfilesService) List(ctx context.Context, actor dto.Actor, role repositories.UserRole) ([]dto.Profile, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateStaffProfileRole(role); err != nil {
		return nil, err
	}

	profiles, err := service.profiles.List(ctx, role)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Profile, 0, len(profiles))
	for _, profile := range profiles {
		output = append(output, *toProfileOutput(&profile))
	}

	return output, nil
}

func (service *staffProfilesService) GetByID(ctx context.Context, actor dto.Actor, role repositories.UserRole, id string) (*dto.Profile, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateStaffProfileRole(role); err != nil {
		return nil, err
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}

	profile, err := service.profiles.GetByID(ctx, id, role)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toProfileOutput(profile), nil
}

func (service *staffProfilesService) Update(ctx context.Context, actor dto.Actor, role repositories.UserRole, id string, input dto.UpdateProfile) (*dto.Profile, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if err := validateStaffProfileRole(role); err != nil {
		return nil, err
	}
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if input.LastName == nil && input.FirstName == nil && input.MiddleName == nil {
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

	profile, err := service.profiles.Update(ctx, id, role, repositories.ProfileUpdateData{
		LastName:      input.LastName,
		FirstName:     input.FirstName,
		MiddleName:    input.MiddleName,
		MiddleNameSet: middleNameSet,
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toProfileOutput(profile), nil
}

func (service *staffProfilesService) Delete(ctx context.Context, actor dto.Actor, role repositories.UserRole, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		return ErrForbidden
	}
	if err := validateStaffProfileRole(role); err != nil {
		return err
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.profiles.Delete(ctx, id, role); err != nil {
		return mapDirectoryRepositoryError(err)
	}

	return nil
}

func validateStaffProfileRole(role repositories.UserRole) error {
	if role != repositories.UserRoleAdministrator && role != repositories.UserRoleTeacher {
		return fmt.Errorf("%w: unsupported staff profile role", ErrInvalidInput)
	}

	return nil
}

func toProfileOutput(profile *repositories.Profile) *dto.Profile {
	return &dto.Profile{
		ID:         profile.ID,
		UserID:     profile.UserID,
		Email:      profile.Email,
		Role:       profile.Role,
		LastName:   profile.LastName,
		FirstName:  profile.FirstName,
		MiddleName: profile.MiddleName,
	}
}
