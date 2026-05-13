package services

import (
	"context"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

func TestGradesServiceCreateAllowsTeacherForOwnCurriculumGroup(t *testing.T) {
	teacherID := testUUID()
	groupID := testUUID()
	curriculumID := testUUID()
	studentID := testUUID()
	value := 5
	grades := newFakeGradeRepository()
	service := NewGradeService(
		grades,
		newFakeStudentRepository(&repositories.Student{ID: studentID, UserID: testUUID(), GroupID: groupID}),
		newFakeCurriculumRepository(&repositories.Curriculum{ID: curriculumID, GroupID: groupID, LeadBy: teacherID}),
		testLogger(),
	)

	grade, err := service.Create(context.Background(), dto.Actor{ID: teacherID, Role: repositories.UserRoleTeacher}, dto.CreateGrade{
		CurriculumID: curriculumID,
		StudentID:    studentID,
		Value:        &value,
	})
	requireNoError(t, err)
	requireEqual(t, teacherID, grade.AuthorID)
	requireEqual(t, teacherID, grades.createData.AuthorID)
}

func TestGradesServiceCreateRejectsTeacherForForeignCurriculum(t *testing.T) {
	groupID := testUUID()
	curriculumID := testUUID()
	studentID := testUUID()
	value := 5
	grades := newFakeGradeRepository()
	service := NewGradeService(
		grades,
		newFakeStudentRepository(&repositories.Student{ID: studentID, UserID: testUUID(), GroupID: groupID}),
		newFakeCurriculumRepository(&repositories.Curriculum{ID: curriculumID, GroupID: groupID, LeadBy: testUUID()}),
		testLogger(),
	)

	_, err := service.Create(context.Background(), testActor(repositories.UserRoleTeacher), dto.CreateGrade{
		CurriculumID: curriculumID,
		StudentID:    studentID,
		Value:        &value,
	})
	requireErrorIs(t, err, ErrForbidden)
	if grades.createCalled {
		t.Fatal("grade repository must not be called for forbidden teacher")
	}
}

func TestGradesServiceCreateRejectsStudentFromAnotherGroup(t *testing.T) {
	teacherID := testUUID()
	curriculumID := testUUID()
	studentID := testUUID()
	value := 5
	grades := newFakeGradeRepository()
	service := NewGradeService(
		grades,
		newFakeStudentRepository(&repositories.Student{ID: studentID, UserID: testUUID(), GroupID: testUUID()}),
		newFakeCurriculumRepository(&repositories.Curriculum{ID: curriculumID, GroupID: testUUID(), LeadBy: teacherID}),
		testLogger(),
	)

	_, err := service.Create(context.Background(), dto.Actor{ID: teacherID, Role: repositories.UserRoleTeacher}, dto.CreateGrade{
		CurriculumID: curriculumID,
		StudentID:    studentID,
		Value:        &value,
	})
	requireErrorIs(t, err, ErrInvalidInput)
	if grades.createCalled {
		t.Fatal("grade repository must not be called when student group does not match curriculum")
	}
}

func TestGradesServiceListReturnsOnlyCurrentStudentGrades(t *testing.T) {
	userID := testUUID()
	studentID := testUUID()
	ownGrade := &repositories.Grade{ID: testUUID(), StudentID: studentID, Value: 5}
	otherGrade := &repositories.Grade{ID: testUUID(), StudentID: testUUID(), Value: 4}
	service := NewGradeService(
		newFakeGradeRepository(ownGrade, otherGrade),
		newFakeStudentRepository(&repositories.Student{ID: studentID, UserID: userID, GroupID: testUUID()}),
		newFakeCurriculumRepository(),
		testLogger(),
	)

	grades, err := service.List(context.Background(), dto.Actor{ID: userID, Role: repositories.UserRoleStudent})
	requireNoError(t, err)
	requireEqual(t, 1, len(grades))
	requireEqual(t, ownGrade.ID, grades[0].ID)
}

func TestGradesServiceDeleteRequiresTeacherAccess(t *testing.T) {
	teacherID := testUUID()
	grade := &repositories.Grade{
		ID:           testUUID(),
		CurriculumID: testUUID(),
		StudentID:    testUUID(),
		Value:        5,
	}
	grades := newFakeGradeRepository(grade)
	service := NewGradeService(
		grades,
		newFakeStudentRepository(&repositories.Student{ID: grade.StudentID, UserID: testUUID(), GroupID: testUUID()}),
		newFakeCurriculumRepository(&repositories.Curriculum{ID: grade.CurriculumID, GroupID: testUUID(), LeadBy: testUUID()}),
		testLogger(),
	)

	err := service.Delete(context.Background(), dto.Actor{ID: teacherID, Role: repositories.UserRoleTeacher}, grade.ID)
	requireErrorIs(t, err, ErrForbidden)
	if grades.deleteCalled {
		t.Fatal("grade repository delete must not be called for forbidden teacher")
	}
}
