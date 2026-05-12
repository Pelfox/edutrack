package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AnalyticsGradeDistribution содержит количество оценок по одному значению.
type AnalyticsGradeDistribution struct {
	Value int
	Count int
}

// AnalyticsSubjectAverage содержит средний балл по одной дисциплине.
type AnalyticsSubjectAverage struct {
	SubjectID    string
	SubjectTitle string
	AverageGrade float64
	GradesCount  int
}

// AnalyticsOverview содержит агрегированные данные для админ-панели.
type AnalyticsOverview struct {
	StudentsCount     int
	TeachersCount     int
	GroupsCount       int
	SpecialtiesCount  int
	SubjectsCount     int
	CurriculumsCount  int
	GradesCount       int
	AverageGrade      *float64
	GradeDistribution []AnalyticsGradeDistribution
	SubjectAverages   []AnalyticsSubjectAverage
}

// AnalyticsRepository описывает методы для чтения агрегированной статистики.
type AnalyticsRepository interface {
	// GetOverview возвращает сводные показатели для админ-панели.
	GetOverview(ctx context.Context) (*AnalyticsOverview, error)
}

// AnalyticsPostgresRepository работает с аналитическими запросами к PostgreSQL.
type AnalyticsPostgresRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

// NewAnalyticsRepository создаёт репозиторий аналитики.
func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsPostgresRepository {
	return &AnalyticsPostgresRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// GetOverview возвращает сводные показатели для админ-панели.
func (repository *AnalyticsPostgresRepository) GetOverview(ctx context.Context) (*AnalyticsOverview, error) {
	overview := &AnalyticsOverview{}

	if err := repository.scanCounts(ctx, overview); err != nil {
		return nil, err
	}
	if err := repository.scanGradeSummary(ctx, overview); err != nil {
		return nil, err
	}

	distribution, err := repository.listGradeDistribution(ctx)
	if err != nil {
		return nil, err
	}
	overview.GradeDistribution = distribution

	averages, err := repository.listSubjectAverages(ctx)
	if err != nil {
		return nil, err
	}
	overview.SubjectAverages = averages

	return overview, nil
}

func (repository *AnalyticsPostgresRepository) scanCounts(ctx context.Context, overview *AnalyticsOverview) error {
	query := `
		SELECT
			(SELECT COUNT(*) FROM students),
			(SELECT COUNT(*) FROM teachers),
			(SELECT COUNT(*) FROM groups),
			(SELECT COUNT(*) FROM specialties),
			(SELECT COUNT(*) FROM subjects),
			(SELECT COUNT(*) FROM curriculums)
	`
	if err := repository.pool.QueryRow(ctx, query).Scan(
		&overview.StudentsCount,
		&overview.TeachersCount,
		&overview.GroupsCount,
		&overview.SpecialtiesCount,
		&overview.SubjectsCount,
		&overview.CurriculumsCount,
	); err != nil {
		return fmt.Errorf("failed to scan analytics counts: %w", err)
	}

	return nil
}

func (repository *AnalyticsPostgresRepository) scanGradeSummary(ctx context.Context, overview *AnalyticsOverview) error {
	query, args, err := repository.builder.
		Select("COUNT(*)", "AVG(value)").
		From("grades").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build grade summary query: %w", err)
	}

	if err := repository.pool.QueryRow(ctx, query, args...).Scan(&overview.GradesCount, &overview.AverageGrade); err != nil {
		return fmt.Errorf("failed to scan grade summary: %w", err)
	}

	return nil
}

func (repository *AnalyticsPostgresRepository) listGradeDistribution(ctx context.Context) ([]AnalyticsGradeDistribution, error) {
	query, args, err := repository.builder.
		Select("value", "COUNT(*)").
		From("grades").
		GroupBy("value").
		OrderBy("value ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build grade distribution query: %w", err)
	}

	rows, err := repository.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list grade distribution: %w", err)
	}
	defer rows.Close()

	items := make([]AnalyticsGradeDistribution, 0)
	for rows.Next() {
		var item AnalyticsGradeDistribution
		if err := rows.Scan(&item.Value, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan grade distribution: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate grade distribution: %w", err)
	}

	return items, nil
}

func (repository *AnalyticsPostgresRepository) listSubjectAverages(ctx context.Context) ([]AnalyticsSubjectAverage, error) {
	query := `
		SELECT subjects.id, subjects.title, AVG(grades.value), COUNT(grades.id)
		FROM subjects
		INNER JOIN curriculums ON curriculums.subject_id = subjects.id
		INNER JOIN grades ON grades.curriculum_id = curriculums.id
		GROUP BY subjects.id, subjects.title
		ORDER BY subjects.title ASC
	`

	rows, err := repository.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subject averages: %w", err)
	}
	defer rows.Close()

	items := make([]AnalyticsSubjectAverage, 0)
	for rows.Next() {
		var item AnalyticsSubjectAverage
		if err := rows.Scan(&item.SubjectID, &item.SubjectTitle, &item.AverageGrade, &item.GradesCount); err != nil {
			return nil, fmt.Errorf("failed to scan subject average: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate subject averages: %w", err)
	}

	return items, nil
}
