package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrDuplicateUserEmail возвращается при нарушении уникальности email пользователя.
var ErrDuplicateUserEmail = errors.New("user email already exists")

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
	Email        *string
	PasswordHash *string
	Role         *UserRole
	UpdatedAt    time.Time
}

// AuthRepository описывает методы, необходимые для авторизации пользователя.
type AuthRepository interface {
	// GetByEmail возвращает пользователя по email.
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// UserRepository описывает методы для работы с пользователями.
type UserRepository interface {
	// Create создаёт пользователя.
	Create(ctx context.Context, data UserCreateData) (*User, error)

	// GetByID возвращает пользователя по идентификатору.
	GetByID(ctx context.Context, id string) (*User, error)

	// Update обновляет пользователя.
	Update(ctx context.Context, id string, data UserUpdateData) (*User, error)

	// Delete удаляет пользователя по идентификатору.
	Delete(ctx context.Context, id string) error
}

// UsersRepository работает с таблицей users.
type UsersRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewUserRepository создаёт репозиторий пользователей.
func NewUserRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт пользователя.
func (repository *UsersRepository) Create(ctx context.Context, data UserCreateData) (*User, error) {
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
func (repository *UsersRepository) GetByID(ctx context.Context, id string) (*User, error) {
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
func (repository *UsersRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
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
func (repository *UsersRepository) Update(ctx context.Context, id string, data UserUpdateData) (*User, error) {
	update := repository.builder.
		Update("users").
		Set("updated_at", data.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, email, password_hash, role, created_at, updated_at")

	if data.Email != nil {
		update = update.Set("email", *data.Email)
	}
	if data.PasswordHash != nil {
		update = update.Set("password_hash", *data.PasswordHash)
	}
	if data.Role != nil {
		update = update.Set("role", *data.Role)
	}

	query, args, err := update.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update user query: %w", err)
	}

	return scanUser(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет пользователя по идентификатору.
func (repository *UsersRepository) Delete(ctx context.Context, id string) error {
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
	var role string
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		if isUniqueViolation(err, "users_email_unique_idx") {
			return nil, ErrDuplicateUserEmail
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	user.Role = UserRole(role)
	return user, nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" && strings.EqualFold(pgErr.ConstraintName, constraint)
}
