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

// Subject содержит поля таблицы subjects.
type Subject struct {
	ID        string
	Title     string
	CreatedAt time.Time
}

// SubjectCreateData содержит данные для создания предмета.
type SubjectCreateData struct {
	ID        string
	Title     string
	CreatedAt time.Time
}

// SubjectUpdateData содержит данные для обновления предмета.
type SubjectUpdateData struct {
	Title string
}

// SubjectRepository описывает методы для работы с предметами.
type SubjectRepository interface {
	// Create создаёт предмет.
	Create(ctx context.Context, data SubjectCreateData) (*Subject, error)

	// List возвращает список предметов.
	List(ctx context.Context) ([]Subject, error)

	// GetByID возвращает предмет по идентификатору.
	GetByID(ctx context.Context, id string) (*Subject, error)

	// Update обновляет предмет.
	Update(ctx context.Context, id string, data SubjectUpdateData) (*Subject, error)

	// Delete удаляет предмет по идентификатору.
	Delete(ctx context.Context, id string) error
}

// SubjectsRepository работает с таблицей subjects.
type SubjectsRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewSubjectRepository создаёт репозиторий предметов.
func NewSubjectRepository(pool *pgxpool.Pool) *SubjectsRepository {
	return &SubjectsRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт предмет.
func (repository *SubjectsRepository) Create(ctx context.Context, data SubjectCreateData) (*Subject, error) {
	query, args, err := repository.builder.
		Insert("subjects").
		Columns("id", "title", "created_at").
		Values(data.ID, data.Title, data.CreatedAt).
		Suffix("RETURNING id, title, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create subject query: %w", err)
	}

	return scanSubject(repository.pool.QueryRow(ctx, query, args...))
}

// List возвращает список предметов.
func (repository *SubjectsRepository) List(ctx context.Context) ([]Subject, error) {
	query, args, err := repository.builder.
		Select("id", "title", "created_at").
		From("subjects").
		OrderBy("title ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list subjects query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list subjects: %w", err)
	}
	defer rows.Close()

	subjects := make([]Subject, 0)
	for rows.Next() {
		subject, err := scanSubject(rows)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, *subject)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate subjects: %w", err)
	}

	return subjects, nil
}

// GetByID возвращает предмет по идентификатору.
func (repository *SubjectsRepository) GetByID(ctx context.Context, id string) (*Subject, error) {
	query, args, err := repository.builder.
		Select("id", "title", "created_at").
		From("subjects").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get subject by id query: %w", err)
	}

	return scanSubject(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет предмет.
func (repository *SubjectsRepository) Update(ctx context.Context, id string, data SubjectUpdateData) (*Subject, error) {
	query, args, err := repository.builder.
		Update("subjects").
		Set("title", data.Title).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, title, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update subject query: %w", err)
	}

	return scanSubject(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет предмет по идентификатору.
func (repository *SubjectsRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("subjects").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete subject query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete subject: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func scanSubject(row pgx.Row) (*Subject, error) {
	subject := &Subject{}
	if err := row.Scan(&subject.ID, &subject.Title, &subject.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan subject: %w", err)
	}

	return subject, nil
}
