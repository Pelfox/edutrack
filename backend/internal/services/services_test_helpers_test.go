package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func testLogger() zerolog.Logger {
	return zerolog.Nop()
}

func testUUID() string {
	return uuid.NewString()
}

func testNow() time.Time {
	return time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
}

func testActor(role repositories.UserRole) dto.Actor {
	return dto.Actor{ID: testUUID(), Role: role}
}

func testPasswordHash(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	return string(hash)
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

func ptr[T any](value T) *T {
	return &value
}

type fakeUserRepository struct {
	user      *repositories.User
	createErr error
	updateErr error
	deleteErr error

	createCalled bool
	updateCalled bool
	deleteCalled bool

	createData repositories.UserCreateData
	updateData repositories.UserUpdateData
}

func (repository *fakeUserRepository) Create(_ context.Context, data repositories.UserCreateData) (*repositories.User, error) {
	repository.createCalled = true
	repository.createData = data
	if repository.createErr != nil {
		return nil, repository.createErr
	}

	user := &repositories.User{
		ID:           data.ID,
		Email:        data.Email,
		PasswordHash: data.PasswordHash,
		Role:         data.Role,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
	repository.user = user
	return user, nil
}

func (repository *fakeUserRepository) GetByID(_ context.Context, id string) (*repositories.User, error) {
	if repository.user == nil || repository.user.ID != id {
		return nil, pgx.ErrNoRows
	}

	return repository.user, nil
}

func (repository *fakeUserRepository) GetByEmail(_ context.Context, email string) (*repositories.User, error) {
	if repository.user == nil || repository.user.Email != email {
		return nil, pgx.ErrNoRows
	}

	return repository.user, nil
}

func (repository *fakeUserRepository) Update(_ context.Context, id string, data repositories.UserUpdateData) (*repositories.User, error) {
	repository.updateCalled = true
	repository.updateData = data
	if repository.updateErr != nil {
		return nil, repository.updateErr
	}
	if repository.user == nil || repository.user.ID != id {
		return nil, pgx.ErrNoRows
	}

	if data.Email != nil {
		repository.user.Email = *data.Email
	}
	if data.PasswordHash != nil {
		repository.user.PasswordHash = *data.PasswordHash
	}
	if data.Role != nil {
		repository.user.Role = *data.Role
	}
	repository.user.UpdatedAt = data.UpdatedAt
	return repository.user, nil
}

func (repository *fakeUserRepository) Delete(_ context.Context, id string) error {
	repository.deleteCalled = true
	if repository.deleteErr != nil {
		return repository.deleteErr
	}
	if repository.user == nil || repository.user.ID != id {
		return pgx.ErrNoRows
	}

	repository.user = nil
	return nil
}

type fakeStudentRepository struct {
	students     map[string]*repositories.Student
	getByUserID  map[string]*repositories.Student
	createErr    error
	updateErr    error
	deleteErr    error
	createCalled bool
	updateCalled bool
	deleteCalled bool
	createData   repositories.StudentCreateData
	updateData   repositories.StudentUpdateData
}

func newFakeStudentRepository(students ...*repositories.Student) *fakeStudentRepository {
	repository := &fakeStudentRepository{
		students:    make(map[string]*repositories.Student),
		getByUserID: make(map[string]*repositories.Student),
	}
	for _, student := range students {
		repository.students[student.ID] = student
		repository.getByUserID[student.UserID] = student
	}

	return repository
}

func (repository *fakeStudentRepository) Create(_ context.Context, data repositories.StudentCreateData) (*repositories.Student, error) {
	repository.createCalled = true
	repository.createData = data
	if repository.createErr != nil {
		return nil, repository.createErr
	}

	student := &repositories.Student{
		ID:         data.ID,
		UserID:     data.UserID,
		GroupID:    data.GroupID,
		LastName:   data.LastName,
		FirstName:  data.FirstName,
		MiddleName: data.MiddleName,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}
	repository.students[student.ID] = student
	repository.getByUserID[student.UserID] = student
	return student, nil
}

func (repository *fakeStudentRepository) List(_ context.Context) ([]repositories.Student, error) {
	students := make([]repositories.Student, 0, len(repository.students))
	for _, student := range repository.students {
		students = append(students, *student)
	}

	return students, nil
}

func (repository *fakeStudentRepository) GetByID(_ context.Context, id string) (*repositories.Student, error) {
	student, ok := repository.students[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return student, nil
}

func (repository *fakeStudentRepository) GetByUserID(_ context.Context, userID string) (*repositories.Student, error) {
	student, ok := repository.getByUserID[userID]
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return student, nil
}

func (repository *fakeStudentRepository) Update(_ context.Context, id string, data repositories.StudentUpdateData) (*repositories.Student, error) {
	repository.updateCalled = true
	repository.updateData = data
	if repository.updateErr != nil {
		return nil, repository.updateErr
	}

	student, ok := repository.students[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	if data.GroupID != nil {
		student.GroupID = *data.GroupID
	}
	if data.LastName != nil {
		student.LastName = *data.LastName
	}
	if data.FirstName != nil {
		student.FirstName = *data.FirstName
	}
	if data.MiddleNameSet {
		student.MiddleName = data.MiddleName
	}
	student.UpdatedAt = data.UpdatedAt
	return student, nil
}

func (repository *fakeStudentRepository) Delete(_ context.Context, id string) error {
	repository.deleteCalled = true
	if repository.deleteErr != nil {
		return repository.deleteErr
	}
	if _, ok := repository.students[id]; !ok {
		return pgx.ErrNoRows
	}
	delete(repository.students, id)
	return nil
}

type fakeCurriculumRepository struct {
	curriculums map[string]*repositories.Curriculum
}

func newFakeCurriculumRepository(curriculums ...*repositories.Curriculum) *fakeCurriculumRepository {
	repository := &fakeCurriculumRepository{curriculums: make(map[string]*repositories.Curriculum)}
	for _, curriculum := range curriculums {
		repository.curriculums[curriculum.ID] = curriculum
	}

	return repository
}

func (repository *fakeCurriculumRepository) Create(_ context.Context, data repositories.CurriculumCreateData) (*repositories.Curriculum, error) {
	curriculum := &repositories.Curriculum{
		ID:         data.ID,
		Hours:      data.Hours,
		Semester:   data.Semester,
		ReportType: data.ReportType,
		SubjectID:  data.SubjectID,
		GroupID:    data.GroupID,
		LeadBy:     data.LeadBy,
	}
	repository.curriculums[curriculum.ID] = curriculum
	return curriculum, nil
}

func (repository *fakeCurriculumRepository) List(_ context.Context) ([]repositories.Curriculum, error) {
	curriculums := make([]repositories.Curriculum, 0, len(repository.curriculums))
	for _, curriculum := range repository.curriculums {
		curriculums = append(curriculums, *curriculum)
	}

	return curriculums, nil
}

func (repository *fakeCurriculumRepository) GetByID(_ context.Context, id string) (*repositories.Curriculum, error) {
	curriculum, ok := repository.curriculums[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return curriculum, nil
}

func (repository *fakeCurriculumRepository) Update(_ context.Context, id string, data repositories.CurriculumUpdateData) (*repositories.Curriculum, error) {
	curriculum, ok := repository.curriculums[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	if data.Hours != nil {
		curriculum.Hours = *data.Hours
	}
	if data.Semester != nil {
		curriculum.Semester = *data.Semester
	}
	if data.ReportType != nil {
		curriculum.ReportType = *data.ReportType
	}
	if data.SubjectID != nil {
		curriculum.SubjectID = *data.SubjectID
	}
	if data.GroupID != nil {
		curriculum.GroupID = *data.GroupID
	}
	if data.LeadBy != nil {
		curriculum.LeadBy = *data.LeadBy
	}
	return curriculum, nil
}

func (repository *fakeCurriculumRepository) Delete(_ context.Context, id string) error {
	if _, ok := repository.curriculums[id]; !ok {
		return pgx.ErrNoRows
	}
	delete(repository.curriculums, id)
	return nil
}

type fakeGradeRepository struct {
	grades       map[string]*repositories.Grade
	createCalled bool
	updateCalled bool
	deleteCalled bool
	createData   repositories.GradeCreateData
	updateData   repositories.GradeUpdateData
}

func newFakeGradeRepository(grades ...*repositories.Grade) *fakeGradeRepository {
	repository := &fakeGradeRepository{grades: make(map[string]*repositories.Grade)}
	for _, grade := range grades {
		repository.grades[grade.ID] = grade
	}

	return repository
}

func (repository *fakeGradeRepository) Create(_ context.Context, data repositories.GradeCreateData) (*repositories.Grade, error) {
	repository.createCalled = true
	repository.createData = data
	grade := &repositories.Grade{
		ID:           data.ID,
		CurriculumID: data.CurriculumID,
		StudentID:    data.StudentID,
		AuthorID:     data.AuthorID,
		Value:        data.Value,
		Comment:      data.Comment,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
	repository.grades[grade.ID] = grade
	return grade, nil
}

func (repository *fakeGradeRepository) List(_ context.Context) ([]repositories.Grade, error) {
	grades := make([]repositories.Grade, 0, len(repository.grades))
	for _, grade := range repository.grades {
		grades = append(grades, *grade)
	}

	return grades, nil
}

func (repository *fakeGradeRepository) ListByStudentID(_ context.Context, studentID string) ([]repositories.Grade, error) {
	grades := make([]repositories.Grade, 0)
	for _, grade := range repository.grades {
		if grade.StudentID == studentID {
			grades = append(grades, *grade)
		}
	}

	return grades, nil
}

func (repository *fakeGradeRepository) GetByID(_ context.Context, id string) (*repositories.Grade, error) {
	grade, ok := repository.grades[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return grade, nil
}

func (repository *fakeGradeRepository) Update(_ context.Context, id string, data repositories.GradeUpdateData) (*repositories.Grade, error) {
	repository.updateCalled = true
	repository.updateData = data
	grade, ok := repository.grades[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	if data.CurriculumID != nil {
		grade.CurriculumID = *data.CurriculumID
	}
	if data.StudentID != nil {
		grade.StudentID = *data.StudentID
	}
	if data.Value != nil {
		grade.Value = *data.Value
	}
	if data.CommentSet {
		grade.Comment = data.Comment
	}
	grade.UpdatedAt = data.UpdatedAt
	return grade, nil
}

func (repository *fakeGradeRepository) Delete(_ context.Context, id string) error {
	repository.deleteCalled = true
	if _, ok := repository.grades[id]; !ok {
		return pgx.ErrNoRows
	}
	delete(repository.grades, id)
	return nil
}

type fakeProfileRepository struct {
	profile      *repositories.Profile
	createCalled bool
	updateCalled bool
	deleteCalled bool
}

func (repository *fakeProfileRepository) Create(_ context.Context, role repositories.UserRole, data repositories.ProfileCreateData) (*repositories.Profile, error) {
	repository.createCalled = true
	repository.profile = &repositories.Profile{
		ID:         data.ID,
		UserID:     data.UserID,
		Role:       role,
		LastName:   data.LastName,
		FirstName:  data.FirstName,
		MiddleName: data.MiddleName,
	}
	return repository.profile, nil
}

func (repository *fakeProfileRepository) List(_ context.Context, role repositories.UserRole) ([]repositories.Profile, error) {
	if repository.profile == nil || repository.profile.Role != role {
		return []repositories.Profile{}, nil
	}
	return []repositories.Profile{*repository.profile}, nil
}

func (repository *fakeProfileRepository) GetByID(_ context.Context, id string, role repositories.UserRole) (*repositories.Profile, error) {
	if repository.profile == nil || repository.profile.ID != id || repository.profile.Role != role {
		return nil, pgx.ErrNoRows
	}
	return repository.profile, nil
}

func (repository *fakeProfileRepository) GetByUserID(_ context.Context, userID string, role repositories.UserRole) (*repositories.Profile, error) {
	if repository.profile == nil || repository.profile.UserID != userID || repository.profile.Role != role {
		return nil, pgx.ErrNoRows
	}
	return repository.profile, nil
}

func (repository *fakeProfileRepository) Update(_ context.Context, id string, role repositories.UserRole, data repositories.ProfileUpdateData) (*repositories.Profile, error) {
	repository.updateCalled = true
	if repository.profile == nil || repository.profile.ID != id || repository.profile.Role != role {
		return nil, pgx.ErrNoRows
	}
	if data.LastName != nil {
		repository.profile.LastName = *data.LastName
	}
	if data.FirstName != nil {
		repository.profile.FirstName = *data.FirstName
	}
	if data.MiddleNameSet {
		repository.profile.MiddleName = data.MiddleName
	}
	return repository.profile, nil
}

func (repository *fakeProfileRepository) Delete(_ context.Context, id string, role repositories.UserRole) error {
	repository.deleteCalled = true
	if repository.profile == nil || repository.profile.ID != id || repository.profile.Role != role {
		return pgx.ErrNoRows
	}
	repository.profile = nil
	return nil
}

type fakeAnalyticsRepository struct {
	overview *repositories.AnalyticsOverview
	called   bool
}

func (repository *fakeAnalyticsRepository) GetOverview(_ context.Context) (*repositories.AnalyticsOverview, error) {
	repository.called = true
	if repository.overview == nil {
		return &repositories.AnalyticsOverview{}, nil
	}
	return repository.overview, nil
}
