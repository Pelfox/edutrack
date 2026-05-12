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

// StudyForm отражает допустимые значения enum study_form в базе данных.
type StudyForm string

const (
	StudyFormFullTime   StudyForm = "full_time"
	StudyFormEvening    StudyForm = "evening"
	StudyFormExtramural StudyForm = "extramural"
)

// Group содержит поля таблицы groups.
type Group struct {
	ID            string
	Name          string
	StudyForm     StudyForm
	AdmissionYear int
	SpecialtyID   string
	CreatedAt     time.Time
}

// GroupCreateData содержит данные для создания группы.
type GroupCreateData struct {
	ID            string
	Name          string
	StudyForm     StudyForm
	AdmissionYear int
	SpecialtyID   string
	CreatedAt     time.Time
}

// GroupUpdateData содержит данные для обновления группы.
type GroupUpdateData struct {
	Name          *string
	StudyForm     *StudyForm
	AdmissionYear *int
	SpecialtyID   *string
}

// GroupRepository описывает методы для работы с группами.
type GroupRepository interface {
	// Create создаёт группу.
	Create(ctx context.Context, data GroupCreateData) (*Group, error)

	// List возвращает список групп.
	List(ctx context.Context) ([]Group, error)

	// GetByID возвращает группу по идентификатору.
	GetByID(ctx context.Context, id string) (*Group, error)

	// Update обновляет группу.
	Update(ctx context.Context, id string, data GroupUpdateData) (*Group, error)

	// Delete удаляет группу по идентификатору.
	Delete(ctx context.Context, id string) error
}

// GroupsRepository работает с таблицей groups.
type GroupsRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewGroupRepository создаёт репозиторий групп.
func NewGroupRepository(pool *pgxpool.Pool) *GroupsRepository {
	return &GroupsRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт группу.
func (repository *GroupsRepository) Create(ctx context.Context, data GroupCreateData) (*Group, error) {
	query, args, err := repository.builder.
		Insert("groups").
		Columns("id", "name", "study_form", "admission_year", "specialty_id", "created_at").
		Values(data.ID, data.Name, data.StudyForm, data.AdmissionYear, data.SpecialtyID, data.CreatedAt).
		Suffix("RETURNING id, name, study_form, admission_year, specialty_id, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create group query: %w", err)
	}

	return scanGroup(repository.pool.QueryRow(ctx, query, args...))
}

// List возвращает список групп.
func (repository *GroupsRepository) List(ctx context.Context) ([]Group, error) {
	query, args, err := repository.builder.
		Select("id", "name", "study_form", "admission_year", "specialty_id", "created_at").
		From("groups").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list groups query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	groups := make([]Group, 0)
	for rows.Next() {
		group, err := scanGroup(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate groups: %w", err)
	}

	return groups, nil
}

// GetByID возвращает группу по идентификатору.
func (repository *GroupsRepository) GetByID(ctx context.Context, id string) (*Group, error) {
	query, args, err := repository.builder.
		Select("id", "name", "study_form", "admission_year", "specialty_id", "created_at").
		From("groups").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get group by id query: %w", err)
	}

	return scanGroup(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет группу.
func (repository *GroupsRepository) Update(ctx context.Context, id string, data GroupUpdateData) (*Group, error) {
	update := repository.builder.
		Update("groups").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, name, study_form, admission_year, specialty_id, created_at")

	if data.Name != nil {
		update = update.Set("name", *data.Name)
	}
	if data.StudyForm != nil {
		update = update.Set("study_form", *data.StudyForm)
	}
	if data.AdmissionYear != nil {
		update = update.Set("admission_year", *data.AdmissionYear)
	}
	if data.SpecialtyID != nil {
		update = update.Set("specialty_id", *data.SpecialtyID)
	}

	query, args, err := update.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update group query: %w", err)
	}

	return scanGroup(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет группу по идентификатору.
func (repository *GroupsRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("groups").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete group query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func scanGroup(row pgx.Row) (*Group, error) {
	group := &Group{}
	var studyForm string
	if err := row.Scan(
		&group.ID,
		&group.Name,
		&studyForm,
		&group.AdmissionYear,
		&group.SpecialtyID,
		&group.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan group: %w", err)
	}

	group.StudyForm = StudyForm(studyForm)
	return group, nil
}
