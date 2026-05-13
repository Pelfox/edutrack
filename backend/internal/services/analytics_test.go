package services

import (
	"context"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

func TestAnalyticsServiceGetOverviewRequiresAdministrator(t *testing.T) {
	repository := &fakeAnalyticsRepository{}
	service := NewAnalyticsService(repository, testLogger())

	_, err := service.GetOverview(context.Background(), testActor(repositories.UserRoleTeacher))
	requireErrorIs(t, err, ErrForbidden)
	if repository.called {
		t.Fatal("analytics repository must not be called for forbidden actor")
	}
}

func TestAnalyticsServiceGetOverviewReturnsData(t *testing.T) {
	repository := &fakeAnalyticsRepository{overview: &repositories.AnalyticsOverview{
		StudentsCount: 2,
		TeachersCount: 1,
		GroupsCount:   1,
		GradesCount:   3,
	}}
	service := NewAnalyticsService(repository, testLogger())

	overview, err := service.GetOverview(context.Background(), testActor(repositories.UserRoleAdministrator))
	requireNoError(t, err)
	requireEqual(t, 2, overview.StudentsCount)
	requireEqual(t, 1, overview.TeachersCount)
	requireEqual(t, 1, overview.GroupsCount)
	requireEqual(t, 3, overview.GradesCount)
}
