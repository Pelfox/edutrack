package repositories

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testPool     *pgxpool.Pool
	testPoolOnce sync.Once
	testPoolErr  error
)

func repositoryTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL to run repository integration tests")
	}
	if err := validateTestDatabaseURL(databaseURL); err != nil {
		t.Fatalf("invalid TEST_DATABASE_URL: %v", err)
	}

	testPoolOnce.Do(func() {
		ctx := context.Background()
		testPool, testPoolErr = pgxpool.New(ctx, databaseURL)
		if testPoolErr != nil {
			return
		}
		testPoolErr = waitTestDatabase(ctx, testPool)
		if testPoolErr != nil {
			return
		}
		testPoolErr = migrateTestDatabase(ctx, testPool)
	})
	if testPoolErr != nil {
		t.Fatalf("failed to prepare test database: %v", testPoolErr)
	}

	resetTestDatabase(t, testPool)
	return testPool
}

func waitTestDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	var err error
	for range 30 {
		if err = pool.Ping(ctx); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}

func validateTestDatabaseURL(databaseURL string) error {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}

	databaseName := strings.TrimPrefix(parsedURL.Path, "/")
	if !strings.Contains(strings.ToLower(databaseName), "test") {
		return errors.New("database name must contain \"test\" to avoid destructive cleanup")
	}

	return nil
}

func migrateTestDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, "DROP SCHEMA public CASCADE"); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, "CREATE SCHEMA public"); err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "migrations", "*.up.sql")
	migrationFiles, err := filepath.Glob(migrationsPath)
	if err != nil {
		return err
	}
	sort.Strings(migrationFiles)

	for _, migrationFile := range migrationFiles {
		query, err := os.ReadFile(migrationFile)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(query)); err != nil {
			return err
		}
	}

	return nil
}

func resetTestDatabase(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		TRUNCATE TABLE
			grades,
			curriculums,
			students,
			teachers,
			administrators,
			"groups",
			subjects,
			specialties,
			users
		RESTART IDENTITY CASCADE
	`
	if _, err := pool.Exec(context.Background(), query); err != nil {
		t.Fatalf("failed to reset test database: %v", err)
	}
}

func testNow() time.Time {
	return time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
}

func testUUID() string {
	return uuid.NewString()
}

func ptr[T any](value T) *T {
	return &value
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}

func requireEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func createTestUser(t *testing.T, pool *pgxpool.Pool, role UserRole, email string) *User {
	t.Helper()

	user, err := NewUserRepository(pool).Create(context.Background(), UserCreateData{
		ID:           testUUID(),
		Email:        email,
		PasswordHash: "password-hash",
		Role:         role,
		CreatedAt:    testNow(),
		UpdatedAt:    testNow(),
	})
	requireNoError(t, err)

	return user
}

func createTestSpecialty(t *testing.T, pool *pgxpool.Pool, title string) *Specialty {
	t.Helper()

	specialty, err := NewSpecialtyRepository(pool).Create(context.Background(), SpecialtyCreateData{
		ID:        testUUID(),
		Title:     title,
		CreatedAt: testNow(),
	})
	requireNoError(t, err)

	return specialty
}

func createTestSubject(t *testing.T, pool *pgxpool.Pool, title string) *Subject {
	t.Helper()

	subject, err := NewSubjectRepository(pool).Create(context.Background(), SubjectCreateData{
		ID:        testUUID(),
		Title:     title,
		CreatedAt: testNow(),
	})
	requireNoError(t, err)

	return subject
}

func createTestGroup(t *testing.T, pool *pgxpool.Pool, specialtyID string, name string) *Group {
	t.Helper()

	group, err := NewGroupRepository(pool).Create(context.Background(), GroupCreateData{
		ID:            testUUID(),
		Name:          name,
		StudyForm:     StudyFormFullTime,
		AdmissionYear: 2026,
		SpecialtyID:   specialtyID,
		CreatedAt:     testNow(),
	})
	requireNoError(t, err)

	return group
}

func createTestStudent(t *testing.T, pool *pgxpool.Pool, userID string, groupID string) *Student {
	t.Helper()

	middleName := "Middle"
	student, err := NewStudentRepository(pool).Create(context.Background(), StudentCreateData{
		ID:         testUUID(),
		UserID:     userID,
		GroupID:    groupID,
		LastName:   "Ivanov",
		FirstName:  "Ivan",
		MiddleName: &middleName,
		CreatedAt:  testNow(),
		UpdatedAt:  testNow(),
	})
	requireNoError(t, err)

	return student
}
