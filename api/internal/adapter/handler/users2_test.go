package handler_test

// ListUsers, GetUser, SendInterest, ListInterests のテスト
// users_test.go の pickup/new/recommend とは別ファイルに分離

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/domain/repository"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// --- mocks ---

type mockListUsersUsecase struct {
	fn func(ctx context.Context, userID uint, input usecase.ListUsersInput) ([]*usecase.RecommendUserResult, error)
}

func (m *mockListUsersUsecase) ListUsers(ctx context.Context, userID uint, input usecase.ListUsersInput) ([]*usecase.RecommendUserResult, error) {
	if m.fn != nil {
		return m.fn(ctx, userID, input)
	}
	return nil, nil
}

type mockGetUserUsecase struct {
	fn func(ctx context.Context, myUserID, targetUserID uint) (*usecase.GetUserResult, error)
}

func (m *mockGetUserUsecase) GetUser(ctx context.Context, myUserID, targetUserID uint) (*usecase.GetUserResult, error) {
	if m.fn != nil {
		return m.fn(ctx, myUserID, targetUserID)
	}
	return nil, nil
}

type mockInterestUsecase2 struct {
	fn func(ctx context.Context, from, to uint) (*usecase.SendInterestResult, error)
}

func (m *mockInterestUsecase2) SendInterest(ctx context.Context, from, to uint) (*usecase.SendInterestResult, error) {
	if m.fn != nil {
		return m.fn(ctx, from, to)
	}
	return &usecase.SendInterestResult{Matched: false}, nil
}

type mockListInterestsUsecase struct {
	fn func(ctx context.Context, userID uint, input usecase.ListInterestsInput) ([]*usecase.RecommendUserResult, error)
}

func (m *mockListInterestsUsecase) ListInterests(ctx context.Context, userID uint, input usecase.ListInterestsInput) ([]*usecase.RecommendUserResult, error) {
	if m.fn != nil {
		return m.fn(ctx, userID, input)
	}
	return nil, nil
}

func newTestEchoUsers2(
	listUsers usecase.ListUsersUsecase,
	getUser usecase.GetUserUsecase,
	interest usecase.InterestUsecase,
	listInterests usecase.ListInterestsUsecase,
) *echo.Echo {
	hub := handler.NewWsHub()
	h := &handler.Handler{
		UsersHandler: handler.NewUsersHandler(
			&mockPickupUsecase{},
			&mockNewUserUsecase{},
			&mockRecommendUsecase{},
			listUsers,
			getUser,
			interest,
			listInterests,
			hub,
		),
	}
	return router.NewRouter(h, hub)
}

func dummyRecommendResult() *usecase.RecommendUserResult {
	bio := "自己紹介"
	return &usecase.RecommendUserResult{
		UserID:          10,
		Name:            "山田 太郎",
		Bio:             &bio,
		CommonCount:     5,
		MatchedThings:   []*repository.ThingRecord{},
		UnmatchedThings: []*repository.ThingRecord{},
	}
}

// --- GET /users ---

func TestListUsers_Success(t *testing.T) {
	uc := &mockListUsersUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListUsersInput) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{dummyRecommendResult()}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(uc, nil, nil, nil), http.MethodGet, "/users", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "山田 太郎")
	assert.Contains(t, body, "users")
}

func TestListUsers_Unauthorized(t *testing.T) {
	rec := do(newTestEchoUsers2(nil, nil, nil, nil), http.MethodGet, "/users", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	body := rec.Body.String()
	assert.Equal(t, 1, strings.Count(body, `"code"`), "401 レスポンスが重複している")
}

func TestListUsers_InternalError(t *testing.T) {
	uc := &mockListUsersUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListUsersInput) ([]*usecase.RecommendUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoUsers2(uc, nil, nil, nil), http.MethodGet, "/users", "", 1)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListUsers_Empty(t *testing.T) {
	uc := &mockListUsersUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListUsersInput) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(uc, nil, nil, nil), http.MethodGet, "/users", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "users")
}

func TestListUsers_NoProfileLeak(t *testing.T) {
	uc := &mockListUsersUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListUsersInput) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{dummyRecommendResult()}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(uc, nil, nil, nil), http.MethodGet, "/users", "", 1)
	body := rec.Body.String()
	assert.NotContains(t, body, "email")
	assert.NotContains(t, body, "password")
}

// --- GET /users/:id ---

func TestGetUser_Success(t *testing.T) {
	uc := &mockGetUserUsecase{
		fn: func(_ context.Context, _, targetID uint) (*usecase.GetUserResult, error) {
			interested := false
			return &usecase.GetUserResult{
				UserID:          targetID,
				Name:            "鈴木 花子",
				IsInterested:    interested,
				MatchedThings:   []*repository.ThingRecord{},
				UnmatchedThings: []*repository.ThingRecord{},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, uc, nil, nil), http.MethodGet, "/users/7", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "鈴木 花子")
	assert.Contains(t, body, "isInterested")
	assert.Contains(t, body, "matchedTags")
	assert.Contains(t, body, "unmatchedTags")
}

func TestGetUser_Unauthorized(t *testing.T) {
	rec := do(newTestEchoUsers2(nil, nil, nil, nil), http.MethodGet, "/users/1", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetUser_NotFound(t *testing.T) {
	uc := &mockGetUserUsecase{
		fn: func(_ context.Context, _, _ uint) (*usecase.GetUserResult, error) {
			return nil, usecase.ErrUserNotFound
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, uc, nil, nil), http.MethodGet, "/users/999", "", 1)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "NOT_FOUND")
}

func TestGetUser_InternalError(t *testing.T) {
	uc := &mockGetUserUsecase{
		fn: func(_ context.Context, _, _ uint) (*usecase.GetUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, uc, nil, nil), http.MethodGet, "/users/1", "", 1)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetUser_NoProfileLeak(t *testing.T) {
	uc := &mockGetUserUsecase{
		fn: func(_ context.Context, _, targetID uint) (*usecase.GetUserResult, error) {
			interested := false
			return &usecase.GetUserResult{
				UserID:          targetID,
				Name:            "テスト",
				IsInterested:    interested,
				MatchedThings:   []*repository.ThingRecord{},
				UnmatchedThings: []*repository.ThingRecord{},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, uc, nil, nil), http.MethodGet, "/users/3", "", 1)
	body := rec.Body.String()
	assert.NotContains(t, body, "email")
	assert.NotContains(t, body, "password")
}

// --- POST /users/:id/interests ---

func TestSendInterest_Success_NoMatch(t *testing.T) {
	uc := &mockInterestUsecase2{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			matched := false
			return &usecase.SendInterestResult{Matched: matched}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, uc, nil), http.MethodPost, "/users/2/interests", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "matched")
	assert.Contains(t, body, "false")
}

func TestSendInterest_Success_Match(t *testing.T) {
	roomID := uint(10)
	uc := &mockInterestUsecase2{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			return &usecase.SendInterestResult{Matched: true, RoomID: &roomID}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, uc, nil), http.MethodPost, "/users/2/interests", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"matched":true`)
	assert.Contains(t, body, `"roomId"`)
}

func TestSendInterest_Unauthorized(t *testing.T) {
	rec := do(newTestEchoUsers2(nil, nil, nil, nil), http.MethodPost, "/users/2/interests", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSendInterest_UserNotFound(t *testing.T) {
	uc := &mockInterestUsecase2{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			return nil, usecase.ErrUserNotFound
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, uc, nil), http.MethodPost, "/users/999/interests", "", 1)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "NOT_FOUND")
}

func TestSendInterest_AlreadyInterested(t *testing.T) {
	uc := &mockInterestUsecase2{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			return nil, usecase.ErrAlreadyInterested
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, uc, nil), http.MethodPost, "/users/2/interests", "", 1)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "ALREADY_INTERESTED")
}

func TestSendInterest_InternalError(t *testing.T) {
	uc := &mockInterestUsecase2{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, uc, nil), http.MethodPost, "/users/2/interests", "", 1)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- GET /interests ---

func TestListInterests_Success(t *testing.T) {
	uc := &mockListInterestsUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListInterestsInput) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{dummyRecommendResult()}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, nil, uc), http.MethodGet, "/interests?direction=sent", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "users")
	assert.Contains(t, rec.Body.String(), "山田 太郎")
}

func TestListInterests_Received(t *testing.T) {
	var capturedDir string
	uc := &mockListInterestsUsecase{
		fn: func(_ context.Context, _ uint, input usecase.ListInterestsInput) ([]*usecase.RecommendUserResult, error) {
			capturedDir = input.Direction
			return []*usecase.RecommendUserResult{}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, nil, uc), http.MethodGet, "/interests?direction=received", "", 1)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "received", capturedDir)
}

func TestListInterests_Unauthorized(t *testing.T) {
	rec := do(newTestEchoUsers2(nil, nil, nil, nil), http.MethodGet, "/interests?direction=sent", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	body := rec.Body.String()
	assert.Equal(t, 1, strings.Count(body, `"code"`), "401 レスポンスが重複している")
}

func TestListInterests_InternalError(t *testing.T) {
	uc := &mockListInterestsUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListInterestsInput) ([]*usecase.RecommendUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, nil, uc), http.MethodGet, "/interests?direction=sent", "", 1)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListInterests_NoProfileLeak(t *testing.T) {
	uc := &mockListInterestsUsecase{
		fn: func(_ context.Context, _ uint, _ usecase.ListInterestsInput) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{dummyRecommendResult()}, nil
		},
	}
	rec := doWithAuth(newTestEchoUsers2(nil, nil, nil, uc), http.MethodGet, "/interests?direction=sent", "", 1)
	body := rec.Body.String()
	assert.NotContains(t, body, "email")
	assert.NotContains(t, body, "password")
}
