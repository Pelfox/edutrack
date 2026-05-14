package controllers

import (
	"net/http"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
)

func TestGradeControllerCreatePassesActorAndRequest(t *testing.T) {
	actor := testActor(repositories.UserRoleTeacher)
	value := 5
	grades := &fakeGradesService{output: &dto.Grade{
		ID:           testUUID(),
		CurriculumID: testUUID(),
		StudentID:    testUUID(),
		AuthorID:     actor.ID,
		Value:        value,
	}}
	router := authenticatedRouter(actor)
	NewGradeController(grades).RegisterRoutes(router.Group("/grades"))

	recorder := requestWithToken(router, http.MethodPost, "/grades", dto.CreateGrade{
		CurriculumID: grades.output.CurriculumID,
		StudentID:    grades.output.StudentID,
		Value:        &value,
	}, "token")

	requireStatus(t, recorder, http.StatusCreated)
	output := decodeResponse[dto.Grade](t, recorder)
	requireEqual(t, grades.output.ID, output.ID)
	requireEqual(t, actor.ID, grades.actor.ID)
	requireEqual(t, value, *grades.createInput.Value)
}

func TestGradeControllerListReturnsGrades(t *testing.T) {
	actor := testActor(repositories.UserRoleStudent)
	grades := &fakeGradesService{listOutput: []dto.Grade{
		{ID: testUUID(), StudentID: testUUID(), Value: 4},
	}}
	router := authenticatedRouter(actor)
	NewGradeController(grades).RegisterRoutes(router.Group("/grades"))

	recorder := requestWithToken(router, http.MethodGet, "/grades", nil, "token")

	requireStatus(t, recorder, http.StatusOK)
	output := decodeResponse[[]dto.Grade](t, recorder)
	requireEqual(t, 1, len(output))
	requireEqual(t, grades.listOutput[0].ID, output[0].ID)
	requireEqual(t, actor.ID, grades.actor.ID)
}

func TestGradeControllerUpdateRejectsInvalidJSON(t *testing.T) {
	actor := testActor(repositories.UserRoleTeacher)
	grades := &fakeGradesService{}
	router := authenticatedRouter(actor)
	NewGradeController(grades).RegisterRoutes(router.Group("/grades"))

	recorder := requestWithToken(router, http.MethodPatch, "/grades/"+testUUID(), map[string]any{"value": "bad"}, "token")

	requireStatus(t, recorder, http.StatusBadRequest)
	if grades.updateCalled {
		t.Fatal("service must not be called for invalid JSON")
	}
}

func TestGradeControllerDeleteMapsForbidden(t *testing.T) {
	actor := testActor(repositories.UserRoleTeacher)
	gradeID := testUUID()
	grades := &fakeGradesService{deleteErr: services.ErrForbidden}
	router := authenticatedRouter(actor)
	NewGradeController(grades).RegisterRoutes(router.Group("/grades"))

	recorder := requestWithToken(router, http.MethodDelete, "/grades/"+gradeID, nil, "token")

	requireStatus(t, recorder, http.StatusForbidden)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "forbidden", output.Error)
	requireEqual(t, gradeID, grades.id)
}
