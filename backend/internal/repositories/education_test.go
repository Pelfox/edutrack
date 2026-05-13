package repositories

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCurriculumsRepositoryCreateListUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewCurriculumRepository(pool)
	ctx := context.Background()
	fixtures := createEducationFixtures(t, pool)

	curriculum, err := repository.Create(ctx, CurriculumCreateData{
		ID:         testUUID(),
		Hours:      120,
		Semester:   1,
		ReportType: CurriculumReportTypeExam,
		SubjectID:  fixtures.subject.ID,
		GroupID:    fixtures.group.ID,
		LeadBy:     fixtures.teacher.ID,
	})
	requireNoError(t, err)
	requireEqual(t, CurriculumReportTypeExam, curriculum.ReportType)

	list, err := repository.List(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, curriculum.ID, list[0].ID)

	updatedHours := 144
	updatedReportType := CurriculumReportTypeDiffTest
	updated, err := repository.Update(ctx, curriculum.ID, CurriculumUpdateData{
		Hours:      &updatedHours,
		ReportType: &updatedReportType,
	})
	requireNoError(t, err)
	requireEqual(t, updatedHours, updated.Hours)
	requireEqual(t, updatedReportType, updated.ReportType)

	requireNoError(t, repository.Delete(ctx, curriculum.ID))
	requireErrorIs(t, repository.Delete(ctx, curriculum.ID), pgx.ErrNoRows)
}

func TestGradesRepositoryCreateListByStudentUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewGradeRepository(pool)
	ctx := context.Background()
	fixtures := createEducationFixtures(t, pool)
	curriculum := createTestCurriculum(t, pool, fixtures)

	comment := "Good work"
	grade, err := repository.Create(ctx, GradeCreateData{
		ID:           testUUID(),
		CurriculumID: curriculum.ID,
		StudentID:    fixtures.student.ID,
		AuthorID:     fixtures.teacher.ID,
		Value:        5,
		Comment:      &comment,
		CreatedAt:    testNow(),
		UpdatedAt:    testNow(),
	})
	requireNoError(t, err)
	requireEqual(t, 5, grade.Value)

	list, err := repository.List(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, grade.ID, list[0].ID)

	studentGrades, err := repository.ListByStudentID(ctx, fixtures.student.ID)
	requireNoError(t, err)
	requireEqual(t, 1, len(studentGrades))
	requireEqual(t, grade.ID, studentGrades[0].ID)

	updatedValue := 4
	updated, err := repository.Update(ctx, grade.ID, GradeUpdateData{
		Value:      &updatedValue,
		CommentSet: true,
		Comment:    nil,
		UpdatedAt:  testNow().AddDate(0, 0, 1),
	})
	requireNoError(t, err)
	requireEqual(t, updatedValue, updated.Value)
	if updated.Comment != nil {
		t.Fatalf("expected comment to be nil, got %q", *updated.Comment)
	}

	requireNoError(t, repository.Delete(ctx, grade.ID))
	requireErrorIs(t, repository.Delete(ctx, grade.ID), pgx.ErrNoRows)
}

func TestAnalyticsRepositoryGetOverview(t *testing.T) {
	pool := repositoryTestPool(t)
	ctx := context.Background()
	fixtures := createEducationFixtures(t, pool)
	curriculum := createTestCurriculum(t, pool, fixtures)
	gradesRepository := NewGradeRepository(pool)

	_, err := gradesRepository.Create(ctx, GradeCreateData{
		ID:           testUUID(),
		CurriculumID: curriculum.ID,
		StudentID:    fixtures.student.ID,
		AuthorID:     fixtures.teacher.ID,
		Value:        5,
		CreatedAt:    testNow(),
		UpdatedAt:    testNow(),
	})
	requireNoError(t, err)
	_, err = gradesRepository.Create(ctx, GradeCreateData{
		ID:           testUUID(),
		CurriculumID: curriculum.ID,
		StudentID:    fixtures.student.ID,
		AuthorID:     fixtures.teacher.ID,
		Value:        3,
		CreatedAt:    testNow().Add(time.Second),
		UpdatedAt:    testNow().Add(time.Second),
	})
	requireNoError(t, err)

	overview, err := NewAnalyticsRepository(pool).GetOverview(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, overview.StudentsCount)
	requireEqual(t, 1, overview.TeachersCount)
	requireEqual(t, 1, overview.GroupsCount)
	requireEqual(t, 1, overview.SpecialtiesCount)
	requireEqual(t, 1, overview.SubjectsCount)
	requireEqual(t, 1, overview.CurriculumsCount)
	requireEqual(t, 2, overview.GradesCount)
	if overview.AverageGrade == nil || math.Abs(*overview.AverageGrade-4) > 0.001 {
		t.Fatalf("expected average grade 4, got %v", overview.AverageGrade)
	}
	requireEqual(t, 2, len(overview.GradeDistribution))
	requireEqual(t, 3, overview.GradeDistribution[0].Value)
	requireEqual(t, 1, overview.GradeDistribution[0].Count)
	requireEqual(t, 5, overview.GradeDistribution[1].Value)
	requireEqual(t, 1, overview.GradeDistribution[1].Count)
	requireEqual(t, 1, len(overview.SubjectAverages))
	requireEqual(t, fixtures.subject.ID, overview.SubjectAverages[0].SubjectID)
	requireEqual(t, 2, overview.SubjectAverages[0].GradesCount)
}

type educationFixtures struct {
	specialty *Specialty
	group     *Group
	subject   *Subject
	teacher   *User
	student   *Student
}

func createEducationFixtures(t *testing.T, pool *pgxpool.Pool) educationFixtures {
	t.Helper()

	specialty := createTestSpecialty(t, pool, "Software Engineering")
	group := createTestGroup(t, pool, specialty.ID, "SE-01")
	subject := createTestSubject(t, pool, "Databases")
	teacher := createTestUser(t, pool, UserRoleTeacher, "teacher@example.com")
	studentUser := createTestUser(t, pool, UserRoleStudent, "student@example.com")
	student := createTestStudent(t, pool, studentUser.ID, group.ID)

	_, err := NewProfileRepository(pool).Create(context.Background(), UserRoleTeacher, ProfileCreateData{
		ID:        testUUID(),
		UserID:    teacher.ID,
		LastName:  "Teacher",
		FirstName: "One",
		CreatedAt: testNow(),
		UpdatedAt: testNow(),
	})
	requireNoError(t, err)

	return educationFixtures{
		specialty: specialty,
		group:     group,
		subject:   subject,
		teacher:   teacher,
		student:   student,
	}
}

func createTestCurriculum(t *testing.T, pool *pgxpool.Pool, fixtures educationFixtures) *Curriculum {
	t.Helper()

	curriculum, err := NewCurriculumRepository(pool).Create(context.Background(), CurriculumCreateData{
		ID:         testUUID(),
		Hours:      120,
		Semester:   1,
		ReportType: CurriculumReportTypeExam,
		SubjectID:  fixtures.subject.ID,
		GroupID:    fixtures.group.ID,
		LeadBy:     fixtures.teacher.ID,
	})
	requireNoError(t, err)

	return curriculum
}
