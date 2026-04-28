package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "star-void-music/backend/db/sqlc"
	"star-void-music/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type mockUserRepo struct {
	getUserByEmailFn func(ctx context.Context, email string) (db.User, error)
	getUsersFn       func(ctx context.Context, limit, offset int32) ([]db.User, error)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, email, passwordHash, role string) (db.User, error) {
	return db.User{}, nil
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return db.User{}, nil
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, email)
	}
	return db.User{}, nil
}

func (m *mockUserRepo) GetUsers(ctx context.Context, limit, offset int32) ([]db.User, error) {
	if m.getUsersFn != nil {
		return m.getUsersFn(ctx, limit, offset)
	}
	return []db.User{}, nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return nil
}

func setupGetUsersRouter(repo service.UserRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUserHandler(service.NewUserService(repo))
	r.GET("/users", h.GetUsers)
	return r
}

func TestGetUsers_WithEmailQuery_ReturnsSingleUser(t *testing.T) {
	expected := db.User{
		ID:        uuid.New(),
		Email:     "john@doe.com",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	repo := &mockUserRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (db.User, error) {
			if email != "john@doe.com" {
				t.Fatalf("unexpected email: %s", email)
			}
			return expected, nil
		},
	}
	r := setupGetUsersRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/users?email=john@doe.com", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body struct {
		Data []db.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Email != expected.Email {
		t.Fatalf("unexpected data: %+v", body.Data)
	}
}

func TestGetUsers_WithEmailQuery_NotFoundReturnsEmptyList(t *testing.T) {
	repo := &mockUserRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (db.User, error) {
			return db.User{}, pgx.ErrNoRows
		},
	}
	r := setupGetUsersRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/users?email=missing@doe.com", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(body.Data) != 0 {
		t.Fatalf("expected empty list, got %+v", body.Data)
	}
}

func TestGetUsers_WithPagination_UsesLimitAndOffset(t *testing.T) {
	expected := []db.User{
		{ID: uuid.New(), Email: "a@doe.com", Role: "user", CreatedAt: time.Now()},
		{ID: uuid.New(), Email: "b@doe.com", Role: "admin", CreatedAt: time.Now()},
	}

	var gotLimit, gotOffset int32
	repo := &mockUserRepo{
		getUsersFn: func(ctx context.Context, limit, offset int32) ([]db.User, error) {
			gotLimit, gotOffset = limit, offset
			return expected, nil
		},
	}
	r := setupGetUsersRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/users?limit=5&offset=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if gotLimit != 5 || gotOffset != 10 {
		t.Fatalf("expected limit=5 offset=10, got limit=%d offset=%d", gotLimit, gotOffset)
	}

	var body struct {
		Data []db.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(body.Data) != 2 {
		t.Fatalf("expected 2 users, got %d", len(body.Data))
	}
}

func TestGetUsers_InvalidLimit_ReturnsBadRequest(t *testing.T) {
	repo := &mockUserRepo{}
	r := setupGetUsersRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/users?limit=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

