package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

// ProfileRepository описывает методы для работы с профилями пользователей.
type ProfileRepository interface {
	// GetByUserID возвращает профиль пользователя по роли и идентификатору пользователя.
	GetByUserID(ctx context.Context, userID string, role UserRole) (*Profile, error)
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
