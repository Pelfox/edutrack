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

// Grade содержит поля таблицы grades.
type Grade struct {
	ID           string
	CurriculumID string
	StudentID    string
	AuthorID     string
	Value        int
	Comment      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GradeCreateData содержит данные для создания оценки.
type GradeCreateData struct {
	ID           string
	CurriculumID string
	StudentID    string
	AuthorID     string
	Value        int
	Comment      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GradeUpdateData содержит данные для обновления оценки.
type GradeUpdateData struct {
	CurriculumID *string
	StudentID    *string
	Value        *int
	Comment      *string
	CommentSet   bool
	UpdatedAt    time.Time
}

// GradeRepository описывает методы для работы с оценками.
type GradeRepository interface {
	// Create создаёт оценку.
	Create(ctx context.Context, data GradeCreateData) (*Grade, error)

	// List возвращает список оценок.
	List(ctx context.Context) ([]Grade, error)

	// ListByStudentID возвращает список оценок студента.
	ListByStudentID(ctx context.Context, studentID string) ([]Grade, error)

	// GetByID возвращает оценку по идентификатору.
	GetByID(ctx context.Context, id string) (*Grade, error)

	// Update обновляет оценку.
	Update(ctx context.Context, id string, data GradeUpdateData) (*Grade, error)

	// Delete удаляет оценку по идентификатору.
	Delete(ctx context.Context, id string) error
}

// GradesRepository работает с таблицей grades.
type GradesRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewGradeRepository создаёт репозиторий оценок.
func NewGradeRepository(pool *pgxpool.Pool) *GradesRepository {
	return &GradesRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт оценку.
func (repository *GradesRepository) Create(ctx context.Context, data GradeCreateData) (*Grade, error) {
	query, args, err := repository.builder.
		Insert("grades").
		Columns("id", "curriculum_id", "student_id", "author_id", "value", "comment", "created_at", "updated_at").
		Values(data.ID, data.CurriculumID, data.StudentID, data.AuthorID, data.Value, data.Comment, data.CreatedAt, data.UpdatedAt).
		Suffix("RETURNING id, curriculum_id, student_id, author_id, value, comment, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create grade query: %w", err)
	}

	return scanGrade(repository.pool.QueryRow(ctx, query, args...))
}

// List возвращает список оценок.
func (repository *GradesRepository) List(ctx context.Context) ([]Grade, error) {
	query, args, err := repository.builder.
		Select("id", "curriculum_id", "student_id", "author_id", "value", "comment", "created_at", "updated_at").
		From("grades").
		OrderBy("created_at DESC", "id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list grades query: %w", err)
	}

	return repository.list(ctx, query, args...)
}

// ListByStudentID возвращает список оценок студента.
func (repository *GradesRepository) ListByStudentID(ctx context.Context, studentID string) ([]Grade, error) {
	query, args, err := repository.builder.
		Select("id", "curriculum_id", "student_id", "author_id", "value", "comment", "created_at", "updated_at").
		From("grades").
		Where(sq.Eq{"student_id": studentID}).
		OrderBy("created_at DESC", "id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list grades by student id query: %w", err)
	}

	return repository.list(ctx, query, args...)
}

// GetByID возвращает оценку по идентификатору.
func (repository *GradesRepository) GetByID(ctx context.Context, id string) (*Grade, error) {
	query, args, err := repository.builder.
		Select("id", "curriculum_id", "student_id", "author_id", "value", "comment", "created_at", "updated_at").
		From("grades").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get grade by id query: %w", err)
	}

	return scanGrade(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет оценку.
func (repository *GradesRepository) Update(ctx context.Context, id string, data GradeUpdateData) (*Grade, error) {
	update := repository.builder.
		Update("grades").
		Set("updated_at", data.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, curriculum_id, student_id, author_id, value, comment, created_at, updated_at")

	if data.CurriculumID != nil {
		update = update.Set("curriculum_id", *data.CurriculumID)
	}
	if data.StudentID != nil {
		update = update.Set("student_id", *data.StudentID)
	}
	if data.Value != nil {
		update = update.Set("value", *data.Value)
	}
	if data.CommentSet {
		update = update.Set("comment", data.Comment)
	}

	query, args, err := update.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update grade query: %w", err)
	}

	return scanGrade(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет оценку по идентификатору.
func (repository *GradesRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("grades").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete grade query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete grade: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (repository *GradesRepository) list(ctx context.Context, query string, args ...any) ([]Grade, error) {
	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list grades: %w", err)
	}
	defer rows.Close()

	grades := make([]Grade, 0)
	for rows.Next() {
		grade, err := scanGrade(rows)
		if err != nil {
			return nil, err
		}
		grades = append(grades, *grade)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate grades: %w", err)
	}

	return grades, nil
}

func scanGrade(row pgx.Row) (*Grade, error) {
	grade := &Grade{}
	var comment sql.NullString
	if err := row.Scan(
		&grade.ID,
		&grade.CurriculumID,
		&grade.StudentID,
		&grade.AuthorID,
		&grade.Value,
		&comment,
		&grade.CreatedAt,
		&grade.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan grade: %w", err)
	}

	if comment.Valid {
		grade.Comment = &comment.String
	}
	return grade, nil
}
