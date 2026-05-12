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

// Student содержит поля таблицы students.
type Student struct {
	ID         string
	UserID     string
	GroupID    string
	LastName   string
	FirstName  string
	MiddleName *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// StudentCreateData содержит данные для создания студента.
type StudentCreateData struct {
	ID         string
	UserID     string
	GroupID    string
	LastName   string
	FirstName  string
	MiddleName *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// StudentUpdateData содержит данные для обновления студента.
type StudentUpdateData struct {
	GroupID       *string
	LastName      *string
	FirstName     *string
	MiddleName    *string
	MiddleNameSet bool
	UpdatedAt     time.Time
}

// StudentRepository описывает методы для работы со студентами.
type StudentRepository interface {
	// Create создаёт студента.
	Create(ctx context.Context, data StudentCreateData) (*Student, error)

	// List возвращает список студентов.
	List(ctx context.Context) ([]Student, error)

	// GetByID возвращает студента по идентификатору.
	GetByID(ctx context.Context, id string) (*Student, error)

	// GetByUserID возвращает студента по идентификатору пользователя.
	GetByUserID(ctx context.Context, userID string) (*Student, error)

	// Update обновляет студента.
	Update(ctx context.Context, id string, data StudentUpdateData) (*Student, error)

	// Delete удаляет студента по идентификатору.
	Delete(ctx context.Context, id string) error
}

// StudentsRepository работает с таблицей students.
type StudentsRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewStudentRepository создаёт репозиторий студентов.
func NewStudentRepository(pool *pgxpool.Pool) *StudentsRepository {
	return &StudentsRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт студента.
func (repository *StudentsRepository) Create(ctx context.Context, data StudentCreateData) (*Student, error) {
	query, args, err := repository.builder.
		Insert("students").
		Columns("id", "user_id", "group_id", "last_name", "first_name", "middle_name", "created_at", "updated_at").
		Values(data.ID, data.UserID, data.GroupID, data.LastName, data.FirstName, data.MiddleName, data.CreatedAt, data.UpdatedAt).
		Suffix("RETURNING id, user_id, group_id, last_name, first_name, middle_name, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create student query: %w", err)
	}

	return scanStudent(repository.pool.QueryRow(ctx, query, args...))
}

// List возвращает список студентов.
func (repository *StudentsRepository) List(ctx context.Context) ([]Student, error) {
	query, args, err := repository.builder.
		Select("id", "user_id", "group_id", "last_name", "first_name", "middle_name", "created_at", "updated_at").
		From("students").
		OrderBy("last_name ASC", "first_name ASC", "middle_name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list students query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()

	students := make([]Student, 0)
	for rows.Next() {
		student, err := scanStudent(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, *student)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate students: %w", err)
	}

	return students, nil
}

// GetByID возвращает студента по идентификатору.
func (repository *StudentsRepository) GetByID(ctx context.Context, id string) (*Student, error) {
	query, args, err := repository.builder.
		Select("id", "user_id", "group_id", "last_name", "first_name", "middle_name", "created_at", "updated_at").
		From("students").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get student by id query: %w", err)
	}

	return scanStudent(repository.pool.QueryRow(ctx, query, args...))
}

// GetByUserID возвращает студента по идентификатору пользователя.
func (repository *StudentsRepository) GetByUserID(ctx context.Context, userID string) (*Student, error) {
	query, args, err := repository.builder.
		Select("id", "user_id", "group_id", "last_name", "first_name", "middle_name", "created_at", "updated_at").
		From("students").
		Where(sq.Eq{"user_id": userID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get student by user id query: %w", err)
	}

	return scanStudent(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет студента.
func (repository *StudentsRepository) Update(ctx context.Context, id string, data StudentUpdateData) (*Student, error) {
	update := repository.builder.
		Update("students").
		Set("updated_at", data.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, user_id, group_id, last_name, first_name, middle_name, created_at, updated_at")

	if data.GroupID != nil {
		update = update.Set("group_id", *data.GroupID)
	}
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
		return nil, fmt.Errorf("failed to build update student query: %w", err)
	}

	return scanStudent(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет студента по идентификатору.
func (repository *StudentsRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("students").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete student query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func scanStudent(row pgx.Row) (*Student, error) {
	student := &Student{}
	var middleName sql.NullString
	if err := row.Scan(
		&student.ID,
		&student.UserID,
		&student.GroupID,
		&student.LastName,
		&student.FirstName,
		&middleName,
		&student.CreatedAt,
		&student.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan student: %w", err)
	}

	if middleName.Valid {
		student.MiddleName = &middleName.String
	}
	return student, nil
}
