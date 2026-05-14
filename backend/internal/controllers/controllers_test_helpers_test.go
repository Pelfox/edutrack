package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func init() {
	gin.SetMode(gin.TestMode)
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

func testRouter() *gin.Engine {
	router := gin.New()
	router.Use(RequestLoggerMiddleware(zerolog.Nop()))
	return router
}

func request(router http.Handler, method string, path string, body any) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(data)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func requestWithToken(router http.Handler, method string, path string, body any, token string) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(data)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func decodeResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()
	var output T
	if err := json.Unmarshal(recorder.Body.Bytes(), &output); err != nil {
		t.Fatalf("failed to decode response: %v\nbody: %s", err, recorder.Body.String())
	}

	return output
}

func requireStatus(t *testing.T, recorder *httptest.ResponseRecorder, status int) {
	t.Helper()
	if recorder.Code != status {
		t.Fatalf("expected status %d, got %d with body %s", status, recorder.Code, recorder.Body.String())
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

type fakeTokenParser struct {
	actor dto.Actor
	err   error
	token string
}

func (parser *fakeTokenParser) ParseToken(tokenValue string) (dto.Actor, error) {
	parser.token = tokenValue
	if parser.err != nil {
		return dto.Actor{}, parser.err
	}
	return parser.actor, nil
}

type fakeAuthService struct {
	loginResult *dto.LoginResult
	loginErr    error
	loginInput  dto.Login
	actor       dto.Actor
	parseErr    error
	token       string
}

func (service *fakeAuthService) Login(_ context.Context, input dto.Login) (*dto.LoginResult, error) {
	service.loginInput = input
	if service.loginErr != nil {
		return nil, service.loginErr
	}
	return service.loginResult, nil
}

func (service *fakeAuthService) ParseToken(tokenValue string) (dto.Actor, error) {
	service.token = tokenValue
	if service.parseErr != nil {
		return dto.Actor{}, service.parseErr
	}
	return service.actor, nil
}

type fakeUsersService struct {
	output       *dto.User
	err          error
	deleteErr    error
	actor        dto.Actor
	id           string
	createInput  dto.CreateUser
	updateInput  dto.UpdateUser
	createCalled bool
	updateCalled bool
	deleteCalled bool
}

func (service *fakeUsersService) Create(_ context.Context, actor dto.Actor, input dto.CreateUser) (*dto.User, error) {
	service.actor = actor
	service.createInput = input
	service.createCalled = true
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeUsersService) GetByID(_ context.Context, actor dto.Actor, id string) (*dto.User, error) {
	service.actor = actor
	service.id = id
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeUsersService) Update(_ context.Context, actor dto.Actor, id string, input dto.UpdateUser) (*dto.User, error) {
	service.actor = actor
	service.id = id
	service.updateInput = input
	service.updateCalled = true
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeUsersService) Delete(_ context.Context, actor dto.Actor, id string) error {
	service.actor = actor
	service.id = id
	service.deleteCalled = true
	if service.deleteErr != nil {
		return service.deleteErr
	}
	return nil
}

type fakeStudentsService struct {
	output       *dto.Student
	listOutput   []dto.Student
	err          error
	deleteErr    error
	actor        dto.Actor
	id           string
	createInput  dto.CreateStudent
	updateInput  dto.UpdateStudent
	createCalled bool
	updateCalled bool
	deleteCalled bool
}

func (service *fakeStudentsService) Create(_ context.Context, actor dto.Actor, input dto.CreateStudent) (*dto.Student, error) {
	service.actor = actor
	service.createInput = input
	service.createCalled = true
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeStudentsService) List(_ context.Context, actor dto.Actor) ([]dto.Student, error) {
	service.actor = actor
	if service.err != nil {
		return nil, service.err
	}
	return service.listOutput, nil
}

func (service *fakeStudentsService) GetByID(_ context.Context, actor dto.Actor, id string) (*dto.Student, error) {
	service.actor = actor
	service.id = id
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeStudentsService) GetMe(_ context.Context, actor dto.Actor) (*dto.Student, error) {
	service.actor = actor
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeStudentsService) Update(_ context.Context, actor dto.Actor, id string, input dto.UpdateStudent) (*dto.Student, error) {
	service.actor = actor
	service.id = id
	service.updateInput = input
	service.updateCalled = true
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeStudentsService) Delete(_ context.Context, actor dto.Actor, id string) error {
	service.actor = actor
	service.id = id
	service.deleteCalled = true
	if service.deleteErr != nil {
		return service.deleteErr
	}
	return nil
}

type fakeGradesService struct {
	output       *dto.Grade
	listOutput   []dto.Grade
	err          error
	deleteErr    error
	actor        dto.Actor
	id           string
	createInput  dto.CreateGrade
	updateInput  dto.UpdateGrade
	createCalled bool
	updateCalled bool
	deleteCalled bool
}

func (service *fakeGradesService) Create(_ context.Context, actor dto.Actor, input dto.CreateGrade) (*dto.Grade, error) {
	service.actor = actor
	service.createInput = input
	service.createCalled = true
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeGradesService) List(_ context.Context, actor dto.Actor) ([]dto.Grade, error) {
	service.actor = actor
	if service.err != nil {
		return nil, service.err
	}
	return service.listOutput, nil
}

func (service *fakeGradesService) GetByID(_ context.Context, actor dto.Actor, id string) (*dto.Grade, error) {
	service.actor = actor
	service.id = id
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeGradesService) Update(_ context.Context, actor dto.Actor, id string, input dto.UpdateGrade) (*dto.Grade, error) {
	service.actor = actor
	service.id = id
	service.updateInput = input
	service.updateCalled = true
	if service.err != nil {
		return nil, service.err
	}
	return service.output, nil
}

func (service *fakeGradesService) Delete(_ context.Context, actor dto.Actor, id string) error {
	service.actor = actor
	service.id = id
	service.deleteCalled = true
	if service.deleteErr != nil {
		return service.deleteErr
	}
	return nil
}
