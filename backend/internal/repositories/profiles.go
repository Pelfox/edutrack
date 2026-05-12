package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Profile содержит поля профильных таблиц пользователей.
type Profile struct {
	ID         string
	UserID     string
	Email      string
	Role       UserRole
	LastName   string
	FirstName  string
	MiddleName *string
}

// ProfileCreateData содержит данные для создания профильной записи пользователя.
type ProfileCreateData struct {
	ID         string
	UserID     string
	LastName   string
	FirstName  string
	MiddleName *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ProfileUpdateData содержит данные для обновления профильной записи пользователя.
type ProfileUpdateData struct {
	LastName      *string
	FirstName     *string
	MiddleName    *string
	MiddleNameSet bool
	UpdatedAt     time.Time
}

// ProfileRepository описывает методы для работы с профилями пользователей.
type ProfileRepository interface {
	// Create создаёт профиль пользователя с указанной ролью.
	Create(ctx context.Context, role UserRole, data ProfileCreateData) (*Profile, error)

	// List возвращает список профилей пользователей с указанной ролью.
	List(ctx context.Context, role UserRole) ([]Profile, error)

	// GetByID возвращает профиль по идентификатору и роли пользователя.
	GetByID(ctx context.Context, id string, role UserRole) (*Profile, error)

	// GetByUserID возвращает профиль пользователя по роли и идентификатору пользователя.
	GetByUserID(ctx context.Context, userID string, role UserRole) (*Profile, error)

	// Update обновляет профиль пользователя с указанной ролью.
	Update(ctx context.Context, id string, role UserRole, data ProfileUpdateData) (*Profile, error)

	// Delete удаляет профиль пользователя с указанной ролью.
	Delete(ctx context.Context, id string, role UserRole) error
}

// ProfilesRepository работает с профильными таблицами пользователей.
type ProfilesRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewProfileRepository создаёт репозиторий профилей пользователей.
func NewProfileRepository(pool *pgxpool.Pool) *ProfilesRepository {
	return &ProfilesRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт профиль пользователя с указанной ролью.
func (repository *ProfilesRepository) Create(ctx context.Context, role UserRole, data ProfileCreateData) (*Profile, error) {
	table, err := profileTableByRole(role)
	if err != nil {
		return nil, err
	}

	query, args, err := repository.builder.
		Insert(table).
		Columns("id", "user_id", "last_name", "first_name", "middle_name", "created_at", "updated_at").
		Values(data.ID, data.UserID, data.LastName, data.FirstName, data.MiddleName, data.CreatedAt, data.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create profile query: %w", err)
	}

	var id string
	if err := repository.pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return repository.GetByID(ctx, id, role)
}

// List возвращает список профилей пользователей с указанной ролью.
func (repository *ProfilesRepository) List(ctx context.Context, role UserRole) ([]Profile, error) {
	table, err := profileTableByRole(role)
	if err != nil {
		return nil, err
	}

	query, args, err := repository.builder.
		Select("profiles.id", "profiles.user_id", "users.email", "users.role", "profiles.last_name", "profiles.first_name", "profiles.middle_name").
		From(table+" AS profiles").
		Join("users ON users.id = profiles.user_id").
		Where(sq.Eq{"users.role": role}).
		OrderBy("profiles.last_name ASC", "profiles.first_name ASC", "profiles.middle_name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list profiles query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}
	defer rows.Close()

	profiles := make([]Profile, 0)
	for rows.Next() {
		profile, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate profiles: %w", err)
	}

	return profiles, nil
}

// GetByID возвращает профиль по идентификатору и роли пользователя.
func (repository *ProfilesRepository) GetByID(ctx context.Context, id string, role UserRole) (*Profile, error) {
	table, err := profileTableByRole(role)
	if err != nil {
		return nil, err
	}

	query, args, err := repository.builder.
		Select("profiles.id", "profiles.user_id", "users.email", "users.role", "profiles.last_name", "profiles.first_name", "profiles.middle_name").
		From(table + " AS profiles").
		Join("users ON users.id = profiles.user_id").
		Where(sq.Eq{"profiles.id": id, "users.role": role}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get profile by id query: %w", err)
	}

	return scanProfile(repository.pool.QueryRow(ctx, query, args...))
}

// GetByUserID возвращает профиль пользователя по роли и идентификатору пользователя.
func (repository *ProfilesRepository) GetByUserID(ctx context.Context, userID string, role UserRole) (*Profile, error) {
	table, err := profileTableByRole(role)
	if err != nil {
		return nil, err
	}

	query, args, err := repository.builder.
		Select("profiles.id", "profiles.user_id", "users.email", "users.role", "profiles.last_name", "profiles.first_name", "profiles.middle_name").
		From(table + " AS profiles").
		Join("users ON users.id = profiles.user_id").
		Where(sq.Eq{"profiles.user_id": userID, "users.role": role}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get profile by user id query: %w", err)
	}

	return scanProfile(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет профиль пользователя с указанной ролью.
func (repository *ProfilesRepository) Update(ctx context.Context, id string, role UserRole, data ProfileUpdateData) (*Profile, error) {
	table, err := profileTableByRole(role)
	if err != nil {
		return nil, err
	}

	update := repository.builder.
		Update(table).
		Set("updated_at", data.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id")

	if data.LastName != nil {
		update = update.Set("last_name", *data.LastName)
	}
	if data.FirstName != nil {
		update = update.Set("first_name", *data.FirstName)
	}
	if data.MiddleNameSet {
		update = update.Set("middle_name", data.MiddleName)
	}

	query, args, err := update.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update profile query: %w", err)
	}

	var updatedID string
	if err := repository.pool.QueryRow(ctx, query, args...).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return repository.GetByID(ctx, updatedID, role)
}

// Delete удаляет профиль пользователя с указанной ролью.
func (repository *ProfilesRepository) Delete(ctx context.Context, id string, role UserRole) error {
	table, err := profileTableByRole(role)
	if err != nil {
		return err
	}

	query, args, err := repository.builder.
		Delete(table).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete profile query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func profileTableByRole(role UserRole) (string, error) {
	switch role {
	case UserRoleAdministrator:
		return "administrators", nil
	case UserRoleTeacher:
		return "teachers", nil
	case UserRoleStudent:
		return "students", nil
	default:
		return "", fmt.Errorf("unsupported user role: %s", role)
	}
}

func scanProfile(row pgx.Row) (*Profile, error) {
	profile := &Profile{}
	var middleName sql.NullString
	if err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Email,
		&profile.Role,
		&profile.LastName,
		&profile.FirstName,
		&middleName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan profile: %w", err)
	}

	if middleName.Valid {
		profile.MiddleName = &middleName.String
	}
	return profile, nil
}
