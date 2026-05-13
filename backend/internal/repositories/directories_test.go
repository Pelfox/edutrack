package repositories

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestSpecialtiesRepositoryCreateListUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewSpecialtyRepository(pool)
	ctx := context.Background()

	specialty, err := repository.Create(ctx, SpecialtyCreateData{
		ID:        testUUID(),
		Title:     "Software Engineering",
		CreatedAt: testNow(),
	})
	requireNoError(t, err)

	list, err := repository.List(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, specialty.ID, list[0].ID)

	updated, err := repository.Update(ctx, specialty.ID, SpecialtyUpdateData{Title: "Applied Informatics"})
	requireNoError(t, err)
	requireEqual(t, "Applied Informatics", updated.Title)

	requireNoError(t, repository.Delete(ctx, specialty.ID))
	requireErrorIs(t, repository.Delete(ctx, specialty.ID), pgx.ErrNoRows)
}

func TestSubjectsRepositoryCreateListUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewSubjectRepository(pool)
	ctx := context.Background()

	subject, err := repository.Create(ctx, SubjectCreateData{
		ID:        testUUID(),
		Title:     "Databases",
		CreatedAt: testNow(),
	})
	requireNoError(t, err)

	list, err := repository.List(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, subject.ID, list[0].ID)

	updated, err := repository.Update(ctx, subject.ID, SubjectUpdateData{Title: "Distributed Databases"})
	requireNoError(t, err)
	requireEqual(t, "Distributed Databases", updated.Title)

	requireNoError(t, repository.Delete(ctx, subject.ID))
	requireErrorIs(t, repository.Delete(ctx, subject.ID), pgx.ErrNoRows)
}

func TestGroupsRepositoryCreateListUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewGroupRepository(pool)
	ctx := context.Background()
	specialty := createTestSpecialty(t, pool, "Computer Science")

	group, err := repository.Create(ctx, GroupCreateData{
		ID:            testUUID(),
		Name:          "CS-01",
		StudyForm:     StudyFormFullTime,
		AdmissionYear: 2026,
		SpecialtyID:   specialty.ID,
		CreatedAt:     testNow(),
	})
	requireNoError(t, err)
	requireEqual(t, StudyFormFullTime, group.StudyForm)

	list, err := repository.List(ctx)
	requireNoError(t, err)
	requireEqual(t, 1, len(list))
	requireEqual(t, group.ID, list[0].ID)

	updatedName := "CS-02"
	updatedStudyForm := StudyFormEvening
	updatedAdmissionYear := 2027
	updated, err := repository.Update(ctx, group.ID, GroupUpdateData{
		Name:          &updatedName,
		StudyForm:     &updatedStudyForm,
		AdmissionYear: &updatedAdmissionYear,
	})
	requireNoError(t, err)
	requireEqual(t, updatedName, updated.Name)
	requireEqual(t, updatedStudyForm, updated.StudyForm)
	requireEqual(t, updatedAdmissionYear, updated.AdmissionYear)
	requireEqual(t, specialty.ID, updated.SpecialtyID)

	requireNoError(t, repository.Delete(ctx, group.ID))
	requireErrorIs(t, repository.Delete(ctx, group.ID), pgx.ErrNoRows)
}

func TestGroupsRepositoryRejectsMissingSpecialty(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewGroupRepository(pool)

	_, err := repository.Create(context.Background(), GroupCreateData{
		ID:            testUUID(),
		Name:          "CS-01",
		StudyForm:     StudyFormFullTime,
		AdmissionYear: 2026,
		SpecialtyID:   testUUID(),
		CreatedAt:     testNow(),
	})
	if err == nil {
		t.Fatal("expected foreign key error, got nil")
	}
}
