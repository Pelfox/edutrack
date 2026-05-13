package services

import (
	"context"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

func TestProfileServiceGetMeUsesActorRole(t *testing.T) {
	actor := testActor(repositories.UserRoleTeacher)
	profiles := &fakeProfileRepository{profile: &repositories.Profile{
		ID:        testUUID(),
		UserID:    actor.ID,
		Email:     "teacher@example.com",
		Role:      repositories.UserRoleTeacher,
		LastName:  "Teacher",
		FirstName: "One",
	}}
	service := NewProfileService(profiles, testLogger())

	profile, err := service.GetMe(context.Background(), actor)
	requireNoError(t, err)
	requireEqual(t, profiles.profile.ID, profile.ID)
	requireEqual(t, repositories.UserRoleTeacher, profile.Role)
}

func TestStaffProfilesServiceCreateRequiresAdministrator(t *testing.T) {
	profiles := &fakeProfileRepository{}
	users := &fakeUserRepository{user: &repositories.User{ID: testUUID(), Role: repositories.UserRoleTeacher}}
	service := NewStaffProfilesService(profiles, users, testLogger())

	_, err := service.Create(context.Background(), testActor(repositories.UserRoleTeacher), repositories.UserRoleTeacher, dto.CreateProfile{
		UserID:    users.user.ID,
		LastName:  "Teacher",
		FirstName: "One",
	})
	requireErrorIs(t, err, ErrForbidden)
	if profiles.createCalled {
		t.Fatal("profile repository must not be called when actor is not administrator")
	}
}

func TestStaffProfilesServiceCreateRejectsRoleMismatch(t *testing.T) {
	profiles := &fakeProfileRepository{}
	users := &fakeUserRepository{user: &repositories.User{ID: testUUID(), Role: repositories.UserRoleStudent}}
	service := NewStaffProfilesService(profiles, users, testLogger())

	_, err := service.Create(context.Background(), testActor(repositories.UserRoleAdministrator), repositories.UserRoleTeacher, dto.CreateProfile{
		UserID:    users.user.ID,
		LastName:  "Teacher",
		FirstName: "One",
	})
	requireErrorIs(t, err, ErrInvalidInput)
	if profiles.createCalled {
		t.Fatal("profile repository must not be called for role mismatch")
	}
}

func TestStaffProfilesServiceUpdateRejectsEmptyPatch(t *testing.T) {
	profiles := &fakeProfileRepository{profile: &repositories.Profile{
		ID:        testUUID(),
		UserID:    testUUID(),
		Role:      repositories.UserRoleTeacher,
		LastName:  "Teacher",
		FirstName: "One",
	}}
	service := NewStaffProfilesService(profiles, &fakeUserRepository{}, testLogger())

	_, err := service.Update(context.Background(), testActor(repositories.UserRoleAdministrator), repositories.UserRoleTeacher, profiles.profile.ID, dto.UpdateProfile{})
	requireErrorIs(t, err, ErrInvalidInput)
	if profiles.updateCalled {
		t.Fatal("profile repository must not be called for empty patch")
	}
}
