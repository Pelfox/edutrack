package controllers

import (
	"net/http"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
)

func TestStudentControllerListReturnsStudents(t *testing.T) {
	actor := testActor(repositories.UserRoleTeacher)
	students := &fakeStudentsService{listOutput: []dto.Student{
		{ID: testUUID(), UserID: testUUID(), GroupID: testUUID(), LastName: "Ivanov", FirstName: "Ivan"},
	}}
	router := authenticatedRouter(actor)
	NewStudentController(students).RegisterRoutes(router.Group("/students"))

	recorder := requestWithToken(router, http.MethodGet, "/students", nil, "token")

	requireStatus(t, recorder, http.StatusOK)
	output := decodeResponse[[]dto.Student](t, recorder)
	requireEqual(t, 1, len(output))
	requireEqual(t, students.listOutput[0].ID, output[0].ID)
	requireEqual(t, actor.ID, students.actor.ID)
}

func TestStudentControllerCreateRejectsInvalidJSON(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	students := &fakeStudentsService{}
	router := authenticatedRouter(actor)
	NewStudentController(students).RegisterRoutes(router.Group("/students"))

	recorder := requestWithToken(router, http.MethodPost, "/students", map[string]any{"user_id": 10}, "token")

	requireStatus(t, recorder, http.StatusBadRequest)
	if students.createCalled {
		t.Fatal("service must not be called for invalid JSON")
	}
}

func TestStudentControllerGetMePassesActor(t *testing.T) {
	actor := testActor(repositories.UserRoleStudent)
	students := &fakeStudentsService{output: &dto.Student{ID: testUUID(), UserID: actor.ID}}
	router := authenticatedRouter(actor)
	NewStudentController(students).RegisterRoutes(router.Group("/students"))

	recorder := requestWithToken(router, http.MethodGet, "/students/me", nil, "token")

	requireStatus(t, recorder, http.StatusOK)
	output := decodeResponse[dto.Student](t, recorder)
	requireEqual(t, students.output.ID, output.ID)
	requireEqual(t, actor.ID, students.actor.ID)
}

func TestStudentControllerDeleteMapsNotFound(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	studentID := testUUID()
	students := &fakeStudentsService{deleteErr: services.ErrNotFound}
	router := authenticatedRouter(actor)
	NewStudentController(students).RegisterRoutes(router.Group("/students"))

	recorder := requestWithToken(router, http.MethodDelete, "/students/"+studentID, nil, "token")

	requireStatus(t, recorder, http.StatusNotFound)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "not_found", output.Error)
	requireEqual(t, studentID, students.id)
}
