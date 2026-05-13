package services

import (
	"context"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

func TestStudentsServiceCreateRequiresAdministrator(t *testing.T) {
	students := newFakeStudentRepository()
	service := NewStudentService(students, newFakeCurriculumRepository(), testLogger())

	_, err := service.Create(context.Background(), testActor(repositories.UserRoleTeacher), dto.CreateStudent{
		UserID:    testUUID(),
		GroupID:   testUUID(),
		LastName:  "Ivanov",
		FirstName: "Ivan",
	})
	requireErrorIs(t, err, ErrForbidden)
	if students.createCalled {
		t.Fatal("repository must not be called when actor is not administrator")
	}
}

func TestStudentsServiceListFiltersTeacherStudentsByCurriculumGroups(t *testing.T) {
	teacherID := testUUID()
	allowedGroupID := testUUID()
	otherGroupID := testUUID()
	allowedStudent := &repositories.Student{
		ID:        testUUID(),
		UserID:    testUUID(),
		GroupID:   allowedGroupID,
		LastName:  "Allowed",
		FirstName: "Student",
	}
	otherStudent := &repositories.Student{
		ID:        testUUID(),
		UserID:    testUUID(),
		GroupID:   otherGroupID,
		LastName:  "Other",
		FirstName: "Student",
	}
	service := NewStudentService(
		newFakeStudentRepository(allowedStudent, otherStudent),
		newFakeCurriculumRepository(&repositories.Curriculum{
			ID:      testUUID(),
			GroupID: allowedGroupID,
			LeadBy:  teacherID,
		}),
		testLogger(),
	)

	students, err := service.List(context.Background(), dto.Actor{ID: teacherID, Role: repositories.UserRoleTeacher})
	requireNoError(t, err)
	requireEqual(t, 1, len(students))
	requireEqual(t, allowedStudent.ID, students[0].ID)
}

func TestStudentsServiceGetByIDRejectsTeacherOutsideGroups(t *testing.T) {
	teacherID := testUUID()
	student := &repositories.Student{
		ID:        testUUID(),
		UserID:    testUUID(),
		GroupID:   testUUID(),
		LastName:  "Ivanov",
		FirstName: "Ivan",
	}
	service := NewStudentService(
		newFakeStudentRepository(student),
		newFakeCurriculumRepository(&repositories.Curriculum{
			ID:      testUUID(),
			GroupID: testUUID(),
			LeadBy:  teacherID,
		}),
		testLogger(),
	)

	_, err := service.GetByID(context.Background(), dto.Actor{ID: teacherID, Role: repositories.UserRoleTeacher}, student.ID)
	requireErrorIs(t, err, ErrForbidden)
}

func TestStudentsServiceGetMeReturnsCurrentStudent(t *testing.T) {
	userID := testUUID()
	student := &repositories.Student{
		ID:        testUUID(),
		UserID:    userID,
		GroupID:   testUUID(),
		LastName:  "Ivanov",
		FirstName: "Ivan",
	}
	service := NewStudentService(newFakeStudentRepository(student), newFakeCurriculumRepository(), testLogger())

	result, err := service.GetMe(context.Background(), dto.Actor{ID: userID, Role: repositories.UserRoleStudent})
	requireNoError(t, err)
	requireEqual(t, student.ID, result.ID)
}

func TestStudentsServiceUpdateRejectsEmptyPatch(t *testing.T) {
	student := &repositories.Student{
		ID:        testUUID(),
		UserID:    testUUID(),
		GroupID:   testUUID(),
		LastName:  "Ivanov",
		FirstName: "Ivan",
	}
	students := newFakeStudentRepository(student)
	service := NewStudentService(students, newFakeCurriculumRepository(), testLogger())

	_, err := service.Update(context.Background(), testActor(repositories.UserRoleAdministrator), student.ID, dto.UpdateStudent{})
	requireErrorIs(t, err, ErrInvalidInput)
	if students.updateCalled {
		t.Fatal("repository must not be called for empty patch")
	}
}
