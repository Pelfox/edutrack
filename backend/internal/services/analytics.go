package services

import (
	"context"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

// AnalyticsService описывает операции модуля аналитики.
type AnalyticsService interface {
	// GetOverview возвращает сводные показатели для администратора.
	GetOverview(ctx context.Context, actor dto.Actor) (*dto.AnalyticsOverview, error)
}

type analyticsService struct {
	analytics repositories.AnalyticsRepository
}

// NewAnalyticsService создаёт сервис аналитики.
func NewAnalyticsService(analytics repositories.AnalyticsRepository) AnalyticsService {
	return &analyticsService{analytics: analytics}
}

func (service *analyticsService) GetOverview(ctx context.Context, actor dto.Actor) (*dto.AnalyticsOverview, error) {
	if actor.Role != repositories.UserRoleAdministrator {
		return nil, ErrForbidden
	}

	overview, err := service.analytics.GetOverview(ctx)
	if err != nil {
		return nil, err
	}

	return toAnalyticsOverviewOutput(overview), nil
}

func toAnalyticsOverviewOutput(overview *repositories.AnalyticsOverview) *dto.AnalyticsOverview {
	distribution := make([]dto.AnalyticsGradeDistribution, 0, len(overview.GradeDistribution))
	for _, item := range overview.GradeDistribution {
		distribution = append(distribution, dto.AnalyticsGradeDistribution{
			Value: item.Value,
			Count: item.Count,
		})
	}

	averages := make([]dto.AnalyticsSubjectAverage, 0, len(overview.SubjectAverages))
	for _, item := range overview.SubjectAverages {
		averages = append(averages, dto.AnalyticsSubjectAverage{
			SubjectID:    item.SubjectID,
			SubjectTitle: item.SubjectTitle,
			AverageGrade: item.AverageGrade,
			GradesCount:  item.GradesCount,
		})
	}

	return &dto.AnalyticsOverview{
		StudentsCount:     overview.StudentsCount,
		TeachersCount:     overview.TeachersCount,
		GroupsCount:       overview.GroupsCount,
		SpecialtiesCount:  overview.SpecialtiesCount,
		SubjectsCount:     overview.SubjectsCount,
		CurriculumsCount:  overview.CurriculumsCount,
		GradesCount:       overview.GradesCount,
		AverageGrade:      overview.AverageGrade,
		GradeDistribution: distribution,
		SubjectAverages:   averages,
	}
}
