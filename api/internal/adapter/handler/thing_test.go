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

type mockThingUsecase struct {
	search func(ctx context.Context, q string) ([]*repository.ThingRecord, error)
	create func(ctx context.Context, name string) (*repository.ThingRecord, error)
}

func (m *mockThingUsecase) Search(ctx context.Context, q string) ([]*repository.ThingRecord, error) {
	if m.search != nil {
		return m.search(ctx, q)
	}
	return nil, nil
}

func (m *mockThingUsecase) Create(ctx context.Context, name string) (*repository.ThingRecord, error) {
	if m.create != nil {
		return m.create(ctx, name)
	}
	return nil, nil
}

func newTestEchoThing(uc usecase.ThingUsecase) *echo.Echo {
	h := &handler.Handler{
		ThingHandler: handler.NewThingHandler(uc),
	}
	return router.NewRouter(h, handler.NewWsHub())
}

// --- GET /things ---

func TestListThings_Success(t *testing.T) {
	now := time.Now()
	uc := &mockThingUsecase{
		search: func(_ context.Context, q string) ([]*repository.ThingRecord, error) {
			return []*repository.ThingRecord{
				{ID: 1, Name: "Go", Aliases: []string{"golang"}, CreatedAt: now},
				{ID: 2, Name: "TypeScript", Aliases: []string{"ts"}, CreatedAt: now},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodGet, "/things", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "Go")
	assert.Contains(t, body, "TypeScript")
	assert.Contains(t, body, "things")
}

func TestListThings_WithQueryParam(t *testing.T) {
	var capturedQ string
	uc := &mockThingUsecase{
		search: func(_ context.Context, q string) ([]*repository.ThingRecord, error) {
			capturedQ = q
			return []*repository.ThingRecord{}, nil
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodGet, "/things?q=go", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "go", capturedQ)
}

func TestListThings_Empty(t *testing.T) {
	uc := &mockThingUsecase{
		search: func(_ context.Context, _ string) ([]*repository.ThingRecord, error) {
			return []*repository.ThingRecord{}, nil
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodGet, "/things", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "things")
}

func TestListThings_InternalError(t *testing.T) {
	uc := &mockThingUsecase{
		search: func(_ context.Context, _ string) ([]*repository.ThingRecord, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodGet, "/things", "", 1)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- POST /things ---

func TestCreateThing_Success(t *testing.T) {
	now := time.Now()
	uc := &mockThingUsecase{
		create: func(_ context.Context, name string) (*repository.ThingRecord, error) {
			return &repository.ThingRecord{ID: 10, Name: name, Aliases: []string{name}, CreatedAt: now}, nil
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodPost, "/things", `{"name":"Rust"}`, 1)
	assert.Equal(t, http.StatusCreated, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "Rust")
	require.Contains(t, body, `"id"`)
}

func TestCreateThing_AlreadyExists(t *testing.T) {
	uc := &mockThingUsecase{
		create: func(_ context.Context, _ string) (*repository.ThingRecord, error) {
			return nil, usecase.ErrThingAlreadyExists
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodPost, "/things", `{"name":"Go"}`, 1)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "THING_ALREADY_EXISTS")
}

func TestCreateThing_InternalError(t *testing.T) {
	uc := &mockThingUsecase{
		create: func(_ context.Context, _ string) (*repository.ThingRecord, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoThing(uc), http.MethodPost, "/things", `{"name":"test"}`, 1)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
