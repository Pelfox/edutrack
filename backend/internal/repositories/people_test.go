package repositories

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestProfilesRepositoryCreateListUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewProfileRepository(pool)
	ctx := context.Background()
	teacherUser := createTestUser(t, pool, UserRoleTeacher, "teacher@example.com")

	middleName := "Petrovich"
	profile, err := repository.Create(ctx, UserRoleTeacher, ProfileCreateData{
		ID:         testUUID(),
		UserID:     teacherUser.ID,
		LastName:   "Petrov",
		FirstName:  "Petr",
		MiddleName: &middleName,
		CreatedAt:  testNow(),
		UpdatedAt:  testNow(),
	})
	requireNoError(t, err)
	requireEqual(t, teacherUser.Email, profile.Email)
	requireEqual(t, UserRoleTeacher, profile.Role)

	profileByUserID, err := repository.GetByUserID(ctx, teacherUser.ID, UserRoleTeacher)
	requireNoError(t, err)
	requireEqual(t, profile.ID, profileByUserID.ID)

	list, err := repository.List(ctx, UserRoleTeacher)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, profile.ID, list[0].ID)

	updatedLastName := "Sidorov"
	updated, err := repository.Update(ctx, profile.ID, UserRoleTeacher, ProfileUpdateData{
		LastName:      &updatedLastName,
		MiddleNameSet: true,
		MiddleName:    nil,
		UpdatedAt:     testNow().AddDate(0, 0, 1),
	})
	requireNoError(t, err)
	requireEqual(t, updatedLastName, updated.LastName)
	if updated.MiddleName != nil {
		t.Fatalf("expected middle name to be nil, got %q", *updated.MiddleName)
	}

	requireNoError(t, repository.Delete(ctx, profile.ID, UserRoleTeacher))
	_, err = repository.GetByID(ctx, profile.ID, UserRoleTeacher)
	requireErrorIs(t, err, pgx.ErrNoRows)
}

func TestStudentsRepositoryCreateGetByUserUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewStudentRepository(pool)
	ctx := context.Background()
	specialty := createTestSpecialty(t, pool, "Software Engineering")
	group := createTestGroup(t, pool, specialty.ID, "SE-01")
	user := createTestUser(t, pool, UserRoleStudent, "student@example.com")

	student := createTestStudent(t, pool, user.ID, group.ID)

	studentByUserID, err := repository.GetByUserID(ctx, user.ID)
	requireNoError(t, err)
	requireEqual(t, student.ID, studentByUserID.ID)

	list, err := repository.List(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, student.ID, list[0].ID)

	updatedFirstName := "Sergey"
	updated, err := repository.Update(ctx, student.ID, StudentUpdateData{
		FirstName:     &updatedFirstName,
		MiddleNameSet: true,
		MiddleName:    nil,
		UpdatedAt:     testNow().AddDate(0, 0, 1),
	})
	requireNoError(t, err)
	requireEqual(t, updatedFirstName, updated.FirstName)
	if updated.MiddleName != nil {
		t.Fatalf("expected middle name to be nil, got %q", *updated.MiddleName)
	}

	requireNoError(t, repository.Delete(ctx, student.ID))
	_, err = repository.GetByID(ctx, student.ID)
	requireErrorIs(t, err, pgx.ErrNoRows)
}

func TestStudentsRepositoryRejectsDuplicateUser(t *testing.T) {
	pool := repositoryTestPool(t)
	specialty := createTestSpecialty(t, pool, "Software Engineering")
	group := createTestGroup(t, pool, specialty.ID, "SE-01")
	user := createTestUser(t, pool, UserRoleStudent, "student@example.com")

	createTestStudent(t, pool, user.ID, group.ID)
	_, err := NewStudentRepository(pool).Create(context.Background(), StudentCreateData{
		ID:        testUUID(),
		UserID:    user.ID,
		GroupID:   group.ID,
		LastName:  "Ivanov",
		FirstName: "Ivan",
		CreatedAt: testNow(),
		UpdatedAt: testNow(),
	})
	if err == nil {
		t.Fatal("expected unique user_id error, got nil")
	}
}
