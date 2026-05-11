package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRole отражает допустимые значения enum user_role в базе данных.
type UserRole string

const (
	UserRoleAdministrator UserRole = "administrator"
	UserRoleTeacher       UserRole = "teacher"
	UserRoleStudent       UserRole = "student"
)

// User содержит поля таблицы users.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserCreateData содержит данные для создания пользователя.
type UserCreateData struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserUpdateData содержит данные для обновления пользователя.
type UserUpdateData struct {
	Email        string
	PasswordHash string
	Role         UserRole
	UpdatedAt    time.Time
}

// UserRepository работает с таблицей users.
type UserRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewUserRepository создаёт репозиторий пользователей.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт пользователя.
func (repository *UserRepository) Create(ctx context.Context, data UserCreateData) (*User, error) {
	query, args, err := repository.builder.
		Insert("users").
		Columns("id", "email", "password_hash", "role", "created_at", "updated_at").
		Values(data.ID, data.Email, data.PasswordHash, data.Role, data.CreatedAt, data.UpdatedAt).
		Suffix("RETURNING id, email, password_hash, role, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create user query: %w", err)
	}

	return scanUser(repository.pool.QueryRow(ctx, query, args...))
}

// GetByID возвращает пользователя по идентификатору.
func (repository *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query, args, err := repository.builder.
		Select("id", "email", "password_hash", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get user by id query: %w", err)
	}

	return scanUser(repository.pool.QueryRow(ctx, query, args...))
}

// GetByEmail возвращает пользователя по email.
func (repository *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query, args, err := repository.builder.
		Select("id", "email", "password_hash", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get user by email query: %w", err)
	}

	return scanUser(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет пользователя.
func (repository *UserRepository) Update(ctx context.Context, id string, data UserUpdateData) (*User, error) {
	query, args, err := repository.builder.
		Update("users").
		Set("email", data.Email).
		Set("password_hash", data.PasswordHash).
		Set("role", data.Role).
		Set("updated_at", data.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, email, password_hash, role, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update user query: %w", err)
	}

	return scanUser(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет пользователя по идентификатору.
func (repository *UserRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete user query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func scanUser(row pgx.Row) (*User, error) {
	user := &User{}
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}
