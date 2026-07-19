package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/domain/repository"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProfileUsecase struct {
	create func(ctx context.Context, userID uint, input usecase.ProfileInput) (*usecase.ProfileResult, error)
	get    func(ctx context.Context, userID uint) (*usecase.ProfileResult, error)
	update func(ctx context.Context, userID uint, input usecase.ProfileInput) (*usecase.ProfileResult, error)
}

func (m *mockProfileUsecase) Create(ctx context.Context, userID uint, input usecase.ProfileInput) (*usecase.ProfileResult, error) {
	if m.create != nil {
		return m.create(ctx, userID, input)
	}
	return nil, nil
}

func (m *mockProfileUsecase) Get(ctx context.Context, userID uint) (*usecase.ProfileResult, error) {
	if m.get != nil {
		return m.get(ctx, userID)
	}
	return nil, nil
}

func (m *mockProfileUsecase) Update(ctx context.Context, userID uint, input usecase.ProfileInput) (*usecase.ProfileResult, error) {
	if m.update != nil {
		return m.update(ctx, userID, input)
	}
	return nil, nil
}

func newTestEchoProfile(uc usecase.ProfileUsecase) *echo.Echo {
	h := &handler.Handler{
		ProfileHandler: handler.NewProfileHandler(uc),
	}
	return router.NewRouter(h, handler.NewWsHub())
}

func dummyProfileResult(userID uint) *usecase.ProfileResult {
	now := time.Now()
	bio := "自己紹介テキスト"
	return &usecase.ProfileResult{
		Profile: &repository.ProfileRecord{
			UserID:    userID,
			Name:      "テスト太郎",
			Bio:       &bio,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Doings: []*repository.ThingRecord{},
		Wants:  []*repository.ThingRecord{},
	}
}

// --- POST /profiles/me ---

func TestCreateMyProfile_Success(t *testing.T) {
	uc := &mockProfileUsecase{
		create: func(_ context.Context, userID uint, _ usecase.ProfileInput) (*usecase.ProfileResult, error) {
			return dummyProfileResult(userID), nil
		},
	}
	rec := doWithAuth(
		newTestEchoProfile(uc),
		http.MethodPost, "/profiles/me",
		`{"name":"テスト太郎","doingThingIds":[],"wantThingIds":[]}`,
		1,
	)
	assert.Equal(t, http.StatusCreated, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "テスト太郎")
	require.Contains(t, body, "userId")
}

func TestCreateMyProfile_Unauthorized(t *testing.T) {
	uc := &mockProfileUsecase{}
	rec := do(newTestEchoProfile(uc), http.MethodPost, "/profiles/me",
		`{"name":"test","doingThingIds":[],"wantThingIds":[]}`)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateMyProfile_AlreadyExists(t *testing.T) {
	uc := &mockProfileUsecase{
		create: func(_ context.Context, _ uint, _ usecase.ProfileInput) (*usecase.ProfileResult, error) {
			return nil, usecase.ErrProfileAlreadyExists
		},
	}
	rec := doWithAuth(
		newTestEchoProfile(uc),
		http.MethodPost, "/profiles/me",
		`{"name":"test","doingThingIds":[],"wantThingIds":[]}`,
		1,
	)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "PROFILE_ALREADY_EXISTS")
}

func TestCreateMyProfile_InternalError(t *testing.T) {
	uc := &mockProfileUsecase{
		create: func(_ context.Context, _ uint, _ usecase.ProfileInput) (*usecase.ProfileResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(
		newTestEchoProfile(uc),
		http.MethodPost, "/profiles/me",
		`{"name":"test","doingThingIds":[],"wantThingIds":[]}`,
		1,
	)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- GET /profiles/me ---

func TestGetMyProfile_Success(t *testing.T) {
	uc := &mockProfileUsecase{
		get: func(_ context.Context, userID uint) (*usecase.ProfileResult, error) {
			return dummyProfileResult(userID), nil
		},
	}
	rec := doWithAuth(newTestEchoProfile(uc), http.MethodGet, "/profiles/me", "", 5)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "テスト太郎")
}

func TestGetMyProfile_Unauthorized(t *testing.T) {
	uc := &mockProfileUsecase{}
	rec := do(newTestEchoProfile(uc), http.MethodGet, "/profiles/me", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetMyProfile_NotFound(t *testing.T) {
	uc := &mockProfileUsecase{
		get: func(_ context.Context, _ uint) (*usecase.ProfileResult, error) {
			return nil, usecase.ErrProfileNotFound
		},
	}
	rec := doWithAuth(newTestEchoProfile(uc), http.MethodGet, "/profiles/me", "", 1)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "PROFILE_NOT_FOUND")
}

func TestGetMyProfile_NoLeakSensitiveFields(t *testing.T) {
	uc := &mockProfileUsecase{
		get: func(_ context.Context, userID uint) (*usecase.ProfileResult, error) {
			return dummyProfileResult(userID), nil
		},
	}
	rec := doWithAuth(newTestEchoProfile(uc), http.MethodGet, "/profiles/me", "", 1)
	body := rec.Body.String()
	assert.NotContains(t, body, "password")
	assert.NotContains(t, body, "email")
}

// --- PUT /profiles/me ---

func TestUpdateMyProfile_Success(t *testing.T) {
	existingResult := dummyProfileResult(1)
	uc := &mockProfileUsecase{
		get: func(_ context.Context, userID uint) (*usecase.ProfileResult, error) {
			return existingResult, nil
		},
		update: func(_ context.Context, userID uint, input usecase.ProfileInput) (*usecase.ProfileResult, error) {
			r := dummyProfileResult(userID)
			r.Profile.Name = input.Name
			return r, nil
		},
	}
	rec := doWithAuth(
		newTestEchoProfile(uc),
		http.MethodPatch, "/profiles/me",
		`{"name":"更新後の名前"}`,
		1,
	)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "更新後の名前")
}

func TestUpdateMyProfile_Unauthorized(t *testing.T) {
	uc := &mockProfileUsecase{}
	rec := do(newTestEchoProfile(uc), http.MethodPatch, "/profiles/me", `{"name":"test"}`)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateMyProfile_ProfileNotFound(t *testing.T) {
	uc := &mockProfileUsecase{
		get: func(_ context.Context, _ uint) (*usecase.ProfileResult, error) {
			return nil, usecase.ErrProfileNotFound
		},
	}
	rec := doWithAuth(
		newTestEchoProfile(uc),
		http.MethodPatch, "/profiles/me",
		`{"name":"test"}`,
		1,
	)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "PROFILE_NOT_FOUND")
}

func TestUpdateMyProfile_InternalError(t *testing.T) {
	existingResult := dummyProfileResult(1)
	uc := &mockProfileUsecase{
		get: func(_ context.Context, userID uint) (*usecase.ProfileResult, error) {
			return existingResult, nil
		},
		update: func(_ context.Context, _ uint, _ usecase.ProfileInput) (*usecase.ProfileResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(
		newTestEchoProfile(uc),
		http.MethodPatch, "/profiles/me",
		`{"name":"test"}`,
		1,
	)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
