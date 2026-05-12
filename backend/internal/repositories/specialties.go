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

// Specialty содержит поля таблицы specialties.
type Specialty struct {
	ID        string
	Title     string
	CreatedAt time.Time
}

// SpecialtyCreateData содержит данные для создания специальности.
type SpecialtyCreateData struct {
	ID        string
	Title     string
	CreatedAt time.Time
}

// SpecialtyUpdateData содержит данные для обновления специальности.
type SpecialtyUpdateData struct {
	Title string
}

// SpecialtyRepository описывает методы для работы со специальностями.
type SpecialtyRepository interface {
	// Create создаёт специальность.
	Create(ctx context.Context, data SpecialtyCreateData) (*Specialty, error)

	// List возвращает список специальностей.
	List(ctx context.Context) ([]Specialty, error)

	// GetByID возвращает специальность по идентификатору.
	GetByID(ctx context.Context, id string) (*Specialty, error)

	// Update обновляет специальность.
	Update(ctx context.Context, id string, data SpecialtyUpdateData) (*Specialty, error)

	// Delete удаляет специальность по идентификатору.
	Delete(ctx context.Context, id string) error
}

// SpecialtiesRepository работает с таблицей specialties.
type SpecialtiesRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewSpecialtyRepository создаёт репозиторий специальностей.
func NewSpecialtyRepository(pool *pgxpool.Pool) *SpecialtiesRepository {
	return &SpecialtiesRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт специальность.
func (repository *SpecialtiesRepository) Create(ctx context.Context, data SpecialtyCreateData) (*Specialty, error) {
	query, args, err := repository.builder.
		Insert("specialties").
		Columns("id", "title", "created_at").
		Values(data.ID, data.Title, data.CreatedAt).
		Suffix("RETURNING id, title, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create specialty query: %w", err)
	}

	return scanSpecialty(repository.pool.QueryRow(ctx, query, args...))
}

// List возвращает список специальностей.
func (repository *SpecialtiesRepository) List(ctx context.Context) ([]Specialty, error) {
	query, args, err := repository.builder.
		Select("id", "title", "created_at").
		From("specialties").
		OrderBy("title ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list specialties query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list specialties: %w", err)
	}
	defer rows.Close()

	specialties := make([]Specialty, 0)
	for rows.Next() {
		specialty, err := scanSpecialty(rows)
		if err != nil {
			return nil, err
		}
		specialties = append(specialties, *specialty)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate specialties: %w", err)
	}

	return specialties, nil
}

// GetByID возвращает специальность по идентификатору.
func (repository *SpecialtiesRepository) GetByID(ctx context.Context, id string) (*Specialty, error) {
	query, args, err := repository.builder.
		Select("id", "title", "created_at").
		From("specialties").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get specialty by id query: %w", err)
	}

	return scanSpecialty(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет специальность.
func (repository *SpecialtiesRepository) Update(ctx context.Context, id string, data SpecialtyUpdateData) (*Specialty, error) {
	query, args, err := repository.builder.
		Update("specialties").
		Set("title", data.Title).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, title, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update specialty query: %w", err)
	}

	return scanSpecialty(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет специальность по идентификатору.
func (repository *SpecialtiesRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("specialties").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete specialty query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete specialty: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func scanSpecialty(row pgx.Row) (*Specialty, error) {
	specialty := &Specialty{}
	if err := row.Scan(&specialty.ID, &specialty.Title, &specialty.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan specialty: %w", err)
	}

	return specialty, nil
}
