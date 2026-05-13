package repositories

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestUsersRepositoryCreateGetUpdateDelete(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewUserRepository(pool)
	ctx := context.Background()

	user, err := repository.Create(ctx, UserCreateData{
		ID:           testUUID(),
		Email:        "admin@example.com",
		PasswordHash: "hash",
		Role:         UserRoleAdministrator,
		CreatedAt:    testNow(),
		UpdatedAt:    testNow(),
	})
	requireNoError(t, err)
	requireEqual(t, "admin@example.com", user.Email)
	requireEqual(t, UserRoleAdministrator, user.Role)

	userByID, err := repository.GetByID(ctx, user.ID)
	requireNoError(t, err)
	requireEqual(t, user.ID, userByID.ID)

	userByEmail, err := repository.GetByEmail(ctx, user.Email)
	requireNoError(t, err)
	requireEqual(t, user.ID, userByEmail.ID)

	updatedEmail := "teacher@example.com"
	updatedHash := "updated-hash"
	updatedRole := UserRoleTeacher
	updated, err := repository.Update(ctx, user.ID, UserUpdateData{
		Email:        &updatedEmail,
		PasswordHash: &updatedHash,
		Role:         &updatedRole,
		UpdatedAt:    testNow().AddDate(0, 0, 1),
	})
	requireNoError(t, err)
	requireEqual(t, updatedEmail, updated.Email)
	requireEqual(t, updatedHash, updated.PasswordHash)
	requireEqual(t, updatedRole, updated.Role)

	requireNoError(t, repository.Delete(ctx, user.ID))
	_, err = repository.GetByID(ctx, user.ID)
	requireErrorIs(t, err, pgx.ErrNoRows)
	requireErrorIs(t, repository.Delete(ctx, user.ID), pgx.ErrNoRows)
}

func TestUsersRepositoryDuplicateEmail(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewUserRepository(pool)
	ctx := context.Background()

	data := UserCreateData{
		ID:           testUUID(),
		Email:        "duplicate@example.com",
		PasswordHash: "hash",
		Role:         UserRoleStudent,
		CreatedAt:    testNow(),
		UpdatedAt:    testNow(),
	}
	_, err := repository.Create(ctx, data)
	requireNoError(t, err)

	data.ID = testUUID()
	_, err = repository.Create(ctx, data)
	requireErrorIs(t, err, ErrDuplicateUserEmail)
}

func TestUsersRepositoryGetByEmailNoRows(t *testing.T) {
	pool := repositoryTestPool(t)
	repository := NewUserRepository(pool)

	_, err := repository.GetByEmail(context.Background(), "missing@example.com")
	requireErrorIs(t, err, pgx.ErrNoRows)
}
