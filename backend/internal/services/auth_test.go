package services

import (
	"context"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

func TestAuthServiceLoginReturnsToken(t *testing.T) {
	user := &repositories.User{
		ID:           testUUID(),
		Email:        "admin@example.com",
		PasswordHash: testPasswordHash(t, "password123"),
		Role:         repositories.UserRoleAdministrator,
		CreatedAt:    testNow(),
		UpdatedAt:    testNow(),
	}
	service := NewAuthService(&fakeUserRepository{user: user}, "secret", testLogger())

	result, err := service.Login(context.Background(), dto.Login{
		Email:    " ADMIN@example.com ",
		Password: "password123",
	})
	requireNoError(t, err)
	requireEqual(t, user.ID, result.User.ID)
	requireEqual(t, user.Email, result.User.Email)
	if result.Token == "" {
		t.Fatal("expected token to be set")
	}
}

func TestAuthServiceLoginRejectsInvalidCredentials(t *testing.T) {
	user := &repositories.User{
		ID:           testUUID(),
		Email:        "admin@example.com",
		PasswordHash: testPasswordHash(t, "password123"),
		Role:         repositories.UserRoleAdministrator,
	}
	service := NewAuthService(&fakeUserRepository{user: user}, "secret", testLogger())

	_, err := service.Login(context.Background(), dto.Login{
		Email:    "admin@example.com",
		Password: "wrong-password",
	})
	requireErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthServiceParseToken(t *testing.T) {
	user := &repositories.User{
		ID:           testUUID(),
		Email:        "teacher@example.com",
		PasswordHash: testPasswordHash(t, "password123"),
		Role:         repositories.UserRoleTeacher,
	}
	service := NewAuthService(&fakeUserRepository{user: user}, "secret", testLogger())

	result, err := service.Login(context.Background(), dto.Login{
		Email:    user.Email,
		Password: "password123",
	})
	requireNoError(t, err)

	actor, err := service.ParseToken(result.Token)
	requireNoError(t, err)
	requireEqual(t, user.ID, actor.ID)
	requireEqual(t, repositories.UserRoleTeacher, actor.Role)
}

func TestAuthServiceRejectsEmptySigningSecret(t *testing.T) {
	user := &repositories.User{
		ID:           testUUID(),
		Email:        "admin@example.com",
		PasswordHash: testPasswordHash(t, "password123"),
		Role:         repositories.UserRoleAdministrator,
	}
	service := NewAuthService(&fakeUserRepository{user: user}, "", testLogger())

	_, err := service.Login(context.Background(), dto.Login{
		Email:    user.Email,
		Password: "password123",
	})
	requireErrorIs(t, err, ErrTokenSigningSecret)
}
