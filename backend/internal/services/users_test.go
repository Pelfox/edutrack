package services

import (
	"context"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

func TestUsersServiceCreateRequiresAdministrator(t *testing.T) {
	repository := &fakeUserRepository{}
	service := NewUserService(repository, testLogger())

	_, err := service.Create(context.Background(), testActor(repositories.UserRoleTeacher), dto.CreateUser{
		Email:    "student@example.com",
		Password: "password123",
		Role:     repositories.UserRoleStudent,
	})
	requireErrorIs(t, err, ErrForbidden)
	if repository.createCalled {
		t.Fatal("repository must not be called when actor is not administrator")
	}
}

func TestUsersServiceCreateNormalizesEmailAndMapsDuplicate(t *testing.T) {
	repository := &fakeUserRepository{createErr: repositories.ErrDuplicateUserEmail}
	service := NewUserService(repository, testLogger())

	_, err := service.Create(context.Background(), testActor(repositories.UserRoleAdministrator), dto.CreateUser{
		Email:    " STUDENT@example.com ",
		Password: "password123",
		Role:     repositories.UserRoleStudent,
	})
	requireErrorIs(t, err, ErrDuplicateUserEmail)
	requireEqual(t, "student@example.com", repository.createData.Email)
}

func TestUsersServiceUpdateAllowsSelfButRejectsRoleChange(t *testing.T) {
	user := &repositories.User{
		ID:           testUUID(),
		Email:        "student@example.com",
		PasswordHash: "hash",
		Role:         repositories.UserRoleStudent,
	}
	repository := &fakeUserRepository{user: user}
	service := NewUserService(repository, testLogger())

	role := repositories.UserRoleTeacher
	_, err := service.Update(context.Background(), dto.Actor{ID: user.ID, Role: repositories.UserRoleStudent}, user.ID, dto.UpdateUser{
		Role: &role,
	})
	requireErrorIs(t, err, ErrForbidden)
	if repository.updateCalled {
		t.Fatal("repository must not be called for forbidden role update")
	}
}

func TestUsersServiceUpdateAllowsAdministratorRoleChange(t *testing.T) {
	user := &repositories.User{
		ID:           testUUID(),
		Email:        "student@example.com",
		PasswordHash: "hash",
		Role:         repositories.UserRoleStudent,
	}
	repository := &fakeUserRepository{user: user}
	service := NewUserService(repository, testLogger())

	role := repositories.UserRoleTeacher
	updated, err := service.Update(context.Background(), testActor(repositories.UserRoleAdministrator), user.ID, dto.UpdateUser{
		Role: &role,
	})
	requireNoError(t, err)
	requireEqual(t, repositories.UserRoleTeacher, updated.Role)
	requireEqual(t, repositories.UserRoleTeacher, *repository.updateData.Role)
}

func TestUsersServiceDeleteRequiresAdministrator(t *testing.T) {
	repository := &fakeUserRepository{user: &repositories.User{ID: testUUID()}}
	service := NewUserService(repository, testLogger())

	err := service.Delete(context.Background(), testActor(repositories.UserRoleStudent), repository.user.ID)
	requireErrorIs(t, err, ErrForbidden)
	if repository.deleteCalled {
		t.Fatal("repository must not be called when delete is forbidden")
	}
}
