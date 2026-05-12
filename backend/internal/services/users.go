package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// UsersService описывает операции пользовательского модуля.
type UsersService interface {
	// Create создаёт пользователя от имени администратора.
	Create(ctx context.Context, actor dto.Actor, input dto.CreateUser) (*dto.User, error)

	// GetByID возвращает пользователя с учётом прав доступа текущего пользователя.
	GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.User, error)

	// Update обновляет пользователя с учётом прав доступа текущего пользователя.
	Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateUser) (*dto.User, error)

	// Delete удаляет пользователя от имени администратора.
	Delete(ctx context.Context, actor dto.Actor, id string) error
}

type usersService struct {
	users     repositories.UserRepository
	validator *validator.Validate
}

// NewUserService создаёт сервис пользовательского модуля.
func NewUserService(users repositories.UserRepository) UsersService {
	return &usersService{
		users:     users,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (service *usersService) Create(ctx context.Context, actor dto.Actor, input dto.CreateUser) (*dto.User, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	input.Email = normalizeEmail(input.Email)
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user, err := service.users.Create(ctx, repositories.UserCreateData{
		ID:           uuid.NewString(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         input.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, mapUserRepositoryError(err)
	}

	return toUserOutput(user), nil
}

func (service *usersService) GetByID(ctx context.Context, actor dto.Actor, id string) (*dto.User, error) {
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if actor.Role != repositories.UserRoleAdministrator && actor.ID != id {
		return nil, ErrForbidden
	}

	user, err := service.users.GetByID(ctx, id)
	if err != nil {
		return nil, mapUserRepositoryError(err)
	}

	return toUserOutput(user), nil
}

func (service *usersService) Update(ctx context.Context, actor dto.Actor, id string, input dto.UpdateUser) (*dto.User, error) {
	if err := validateUUID(id); err != nil {
		return nil, err
	}
	if actor.Role != repositories.UserRoleAdministrator && actor.ID != id {
		return nil, ErrForbidden
	}
	if input.Role != nil && actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}
	if input.Email == nil && input.Password == nil && input.Role == nil {
		return nil, fmt.Errorf("%w: at least one field is required", ErrInvalidInput)
	}
	if input.Email != nil {
		email := normalizeEmail(*input.Email)
		input.Email = &email
	}
	if err := validateStruct(service.validator, input); err != nil {
		return nil, err
	}

	data := repositories.UserUpdateData{
		UpdatedAt: time.Now().UTC(),
	}

	if input.Email != nil {
		data.Email = input.Email
	}
	if input.Password != nil {
		passwordHash, err := hashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		data.PasswordHash = &passwordHash
	}
	if input.Role != nil {
		data.Role = input.Role
	}

	user, err := service.users.Update(ctx, id, data)
	if err != nil {
		return nil, mapUserRepositoryError(err)
	}

	return toUserOutput(user), nil
}

func (service *usersService) Delete(ctx context.Context, actor dto.Actor, id string) error {
	if actor.Role != repositories.UserRoleAdministrator {
		return ErrForbidden
	}
	if err := validateUUID(id); err != nil {
		return err
	}

	if err := service.users.Delete(ctx, id); err != nil {
		return mapUserRepositoryError(err)
	}

	return nil
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validateStruct(validate *validator.Validate, value any) error {
	if err := validate.Struct(value); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	return nil
}

func validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: id is invalid", ErrInvalidInput)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func mapUserRepositoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	if errors.Is(err, repositories.ErrDuplicateUserEmail) {
		return ErrDuplicateUserEmail
	}

	return err
}

func toUserOutput(user *repositories.User) *dto.User {
	return &dto.User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
