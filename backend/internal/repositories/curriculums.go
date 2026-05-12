package repositories

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CurriculumReportType отражает допустимые значения enum curriculum_report_type в базе данных.
type CurriculumReportType string

const (
	CurriculumReportTypeExam     CurriculumReportType = "exam"
	CurriculumReportTypeTest     CurriculumReportType = "test"
	CurriculumReportTypeDiffTest CurriculumReportType = "diff_test"
)

// Curriculum содержит поля таблицы curriculums.
type Curriculum struct {
	ID         string
	Hours      int
	Semester   int
	ReportType CurriculumReportType
	SubjectID  string
	GroupID    string
	LeadBy     string
}

// CurriculumCreateData содержит данные для создания учебного плана.
type CurriculumCreateData struct {
	ID         string
	Hours      int
	Semester   int
	ReportType CurriculumReportType
	SubjectID  string
	GroupID    string
	LeadBy     string
}

// CurriculumUpdateData содержит данные для обновления учебного плана.
type CurriculumUpdateData struct {
	Hours      *int
	Semester   *int
	ReportType *CurriculumReportType
	SubjectID  *string
	GroupID    *string
	LeadBy     *string
}

// CurriculumRepository описывает методы для работы с учебными планами.
type CurriculumRepository interface {
	// Create создаёт учебный план.
	Create(ctx context.Context, data CurriculumCreateData) (*Curriculum, error)

	// List возвращает список учебных планов.
	List(ctx context.Context) ([]Curriculum, error)

	// GetByID возвращает учебный план по идентификатору.
	GetByID(ctx context.Context, id string) (*Curriculum, error)

	// Update обновляет учебный план.
	Update(ctx context.Context, id string, data CurriculumUpdateData) (*Curriculum, error)

	// Delete удаляет учебный план по идентификатору.
	Delete(ctx context.Context, id string) error
}

// CurriculumsRepository работает с таблицей curriculums.
type CurriculumsRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewCurriculumRepository создаёт репозиторий учебных планов.
func NewCurriculumRepository(pool *pgxpool.Pool) *CurriculumsRepository {
	return &CurriculumsRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create создаёт учебный план.
func (repository *CurriculumsRepository) Create(ctx context.Context, data CurriculumCreateData) (*Curriculum, error) {
	query, args, err := repository.builder.
		Insert("curriculums").
		Columns("id", "hours", "semester", "report_type", "subject_id", "group_id", "lead_by").
		Values(data.ID, data.Hours, data.Semester, data.ReportType, data.SubjectID, data.GroupID, data.LeadBy).
		Suffix("RETURNING id, hours, semester, report_type, subject_id, group_id, lead_by").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create curriculum query: %w", err)
	}

	return scanCurriculum(repository.pool.QueryRow(ctx, query, args...))
}

// List возвращает список учебных планов.
func (repository *CurriculumsRepository) List(ctx context.Context) ([]Curriculum, error) {
	query, args, err := repository.builder.
		Select("id", "hours", "semester", "report_type", "subject_id", "group_id", "lead_by").
		From("curriculums").
		OrderBy("semester ASC", "id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list curriculums query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list curriculums: %w", err)
	}
	defer rows.Close()

	curriculums := make([]Curriculum, 0)
	for rows.Next() {
		curriculum, err := scanCurriculum(rows)
		if err != nil {
			return nil, err
		}
		curriculums = append(curriculums, *curriculum)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate curriculums: %w", err)
	}

	return curriculums, nil
}

// GetByID возвращает учебный план по идентификатору.
func (repository *CurriculumsRepository) GetByID(ctx context.Context, id string) (*Curriculum, error) {
	query, args, err := repository.builder.
		Select("id", "hours", "semester", "report_type", "subject_id", "group_id", "lead_by").
		From("curriculums").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get curriculum by id query: %w", err)
	}

	return scanCurriculum(repository.pool.QueryRow(ctx, query, args...))
}

// Update обновляет учебный план.
func (repository *CurriculumsRepository) Update(ctx context.Context, id string, data CurriculumUpdateData) (*Curriculum, error) {
	update := repository.builder.
		Update("curriculums").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, hours, semester, report_type, subject_id, group_id, lead_by")

	if data.Hours != nil {
		update = update.Set("hours", *data.Hours)
	}
	if data.Semester != nil {
		update = update.Set("semester", *data.Semester)
	}
	if data.ReportType != nil {
		update = update.Set("report_type", *data.ReportType)
	}
	if data.SubjectID != nil {
		update = update.Set("subject_id", *data.SubjectID)
	}
	if data.GroupID != nil {
		update = update.Set("group_id", *data.GroupID)
	}
	if data.LeadBy != nil {
		update = update.Set("lead_by", *data.LeadBy)
	}

	query, args, err := update.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update curriculum query: %w", err)
	}

	return scanCurriculum(repository.pool.QueryRow(ctx, query, args...))
}

// Delete удаляет учебный план по идентификатору.
func (repository *CurriculumsRepository) Delete(ctx context.Context, id string) error {
	query, args, err := repository.builder.
		Delete("curriculums").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete curriculum query: %w", err)
	}

	result, err := repository.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete curriculum: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func scanCurriculum(row pgx.Row) (*Curriculum, error) {
	curriculum := &Curriculum{}
	var reportType string
	if err := row.Scan(
		&curriculum.ID,
		&curriculum.Hours,
		&curriculum.Semester,
		&reportType,
		&curriculum.SubjectID,
		&curriculum.GroupID,
		&curriculum.LeadBy,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan curriculum: %w", err)
	}

	curriculum.ReportType = CurriculumReportType(reportType)
	return curriculum, nil
}
