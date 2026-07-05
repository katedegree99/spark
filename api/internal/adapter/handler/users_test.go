package handler_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mocks ---

type mockPickupUsecase struct {
	listPickup func(ctx context.Context, userID uint) ([]*usecase.PickupUserResult, error)
}

func (m *mockPickupUsecase) ListPickup(ctx context.Context, userID uint) ([]*usecase.PickupUserResult, error) {
	if m.listPickup != nil {
		return m.listPickup(ctx, userID)
	}
	return nil, nil
}

type mockNewUserUsecase struct {
	listNew func(ctx context.Context, userID uint) ([]*usecase.NewUserResult, error)
}

func (m *mockNewUserUsecase) ListNew(ctx context.Context, userID uint) ([]*usecase.NewUserResult, error) {
	if m.listNew != nil {
		return m.listNew(ctx, userID)
	}
	return nil, nil
}

type mockRecommendUsecase struct {
	listRecommend func(ctx context.Context, userID uint) ([]*usecase.RecommendUserResult, error)
}

func (m *mockRecommendUsecase) ListRecommend(ctx context.Context, userID uint) ([]*usecase.RecommendUserResult, error) {
	if m.listRecommend != nil {
		return m.listRecommend(ctx, userID)
	}
	return nil, nil
}

func newTestEchoWithUsers(pickup usecase.PickupUsecase, newUser usecase.NewUserUsecase, recommend usecase.RecommendUsecase) *echo.Echo {
	h := &handler.Handler{
		UsersHandler: handler.NewUsersHandler(pickup, newUser, recommend, nil, nil, nil),
	}
	return router.NewRouter(h)
}

// --- GET /users/pickup ---

func TestListPickupUsers_Success(t *testing.T) {
	bio := "自己紹介"
	iconURL := "https://example.com/icon.jpg"
	uc := &mockPickupUsecase{
		listPickup: func(_ context.Context, _ uint) ([]*usecase.PickupUserResult, error) {
			return []*usecase.PickupUserResult{
				{
					UserID:          42,
					Name:            "堺 理人",
					Bio:             &bio,
					IconURL:         &iconURL,
					MatchedThings:   nil,
					UnmatchedThings: nil,
				},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(uc, &mockNewUserUsecase{}, &mockRecommendUsecase{}), http.MethodGet, "/users/pickup", "", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "堺 理人")
	require.Contains(t, body, "matchedTags")
	require.Contains(t, body, "unmatchedTags")
}

func TestListPickupUsers_Unauthorized(t *testing.T) {
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, &mockNewUserUsecase{}, &mockRecommendUsecase{}), http.MethodGet, "/users/pickup", "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	// 401 レスポンスが一つだけであること（二重レスポンスのリグレッション防止）
	body := rec.Body.String()
	assert.Equal(t, 1, strings.Count(body, `"code"`), "401 レスポンスが重複している")
}

func TestListPickupUsers_InternalError(t *testing.T) {
	uc := &mockPickupUsecase{
		listPickup: func(_ context.Context, _ uint) ([]*usecase.PickupUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(uc, &mockNewUserUsecase{}, &mockRecommendUsecase{}), http.MethodGet, "/users/pickup", "", 1)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListPickupUsers_Empty(t *testing.T) {
	uc := &mockPickupUsecase{
		listPickup: func(_ context.Context, _ uint) ([]*usecase.PickupUserResult, error) {
			return []*usecase.PickupUserResult{}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(uc, &mockNewUserUsecase{}, &mockRecommendUsecase{}), http.MethodGet, "/users/pickup", "", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "users")
}

func TestListPickupUsers_NoProfileLeak(t *testing.T) {
	email := "secret@example.com"
	_ = email
	uc := &mockPickupUsecase{
		listPickup: func(_ context.Context, _ uint) ([]*usecase.PickupUserResult, error) {
			return []*usecase.PickupUserResult{
				{UserID: 1, Name: "田中"},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(uc, &mockNewUserUsecase{}, &mockRecommendUsecase{}), http.MethodGet, "/users/pickup", "", 1)

	body := rec.Body.String()
	assert.NotContains(t, body, "email")
	assert.NotContains(t, body, "password")
}

// --- GET /users/new ---

func TestListNewUsers_Success(t *testing.T) {
	iconURL := "https://example.com/icon.jpg"
	uc := &mockNewUserUsecase{
		listNew: func(_ context.Context, _ uint) ([]*usecase.NewUserResult, error) {
			return []*usecase.NewUserResult{
				{UserID: 10, Name: "酒井 悠人", IconURL: &iconURL},
				{UserID: 11, Name: "田中 花子", IconURL: nil},
			}, nil
		},
	}
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, uc, &mockRecommendUsecase{}), http.MethodGet, "/users/new", "")

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "酒井 悠人")
	assert.Contains(t, body, "田中 花子")
	assert.Contains(t, body, "userId")
	assert.Contains(t, body, "iconUrl")
}

func TestListNewUsers_NoAuthRequired(t *testing.T) {
	uc := &mockNewUserUsecase{
		listNew: func(_ context.Context, _ uint) ([]*usecase.NewUserResult, error) {
			return []*usecase.NewUserResult{}, nil
		},
	}
	// トークンなしでも 200 が返ること
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, uc, &mockRecommendUsecase{}), http.MethodGet, "/users/new", "")

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestListNewUsers_InternalError(t *testing.T) {
	uc := &mockNewUserUsecase{
		listNew: func(_ context.Context, _ uint) ([]*usecase.NewUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, uc, &mockRecommendUsecase{}), http.MethodGet, "/users/new", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListNewUsers_NoProfileLeak(t *testing.T) {
	uc := &mockNewUserUsecase{
		listNew: func(_ context.Context, _ uint) ([]*usecase.NewUserResult, error) {
			return []*usecase.NewUserResult{
				{UserID: 1, Name: "田中"},
			}, nil
		},
	}
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, uc, &mockRecommendUsecase{}), http.MethodGet, "/users/new", "")

	body := rec.Body.String()
	assert.NotContains(t, body, "email")
	assert.NotContains(t, body, "password")
	assert.NotContains(t, body, "bio")
}

func TestListNewUsers_UserIdAndNameAlwaysPresent(t *testing.T) {
	uc := &mockNewUserUsecase{
		listNew: func(_ context.Context, _ uint) ([]*usecase.NewUserResult, error) {
			return []*usecase.NewUserResult{
				{UserID: 5, Name: "山田 太郎"},
			}, nil
		},
	}
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, uc, &mockRecommendUsecase{}), http.MethodGet, "/users/new", "")

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"userId"`)
	assert.Contains(t, body, `"name"`)
}

// --- GET /users/recommend ---

func TestListRecommendUsers_Success(t *testing.T) {
	bio := "自己紹介"
	iconURL := "https://example.com/icon.jpg"
	uc := &mockRecommendUsecase{
		listRecommend: func(_ context.Context, _ uint) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{
				{
					UserID:      42,
					Name:        "小林 田穂",
					Bio:         &bio,
					IconURL:     &iconURL,
					CommonCount: 40,
				},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(&mockPickupUsecase{}, &mockNewUserUsecase{}, uc), http.MethodGet, "/users/recommend", "", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "小林 田穂")
	assert.Contains(t, body, "commonCount")
	assert.Contains(t, body, "matchedTags")
	assert.Contains(t, body, "unmatchedTags")
}

func TestListRecommendUsers_Unauthorized(t *testing.T) {
	rec := do(newTestEchoWithUsers(&mockPickupUsecase{}, &mockNewUserUsecase{}, &mockRecommendUsecase{}), http.MethodGet, "/users/recommend", "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	body := rec.Body.String()
	assert.Equal(t, 1, strings.Count(body, `"code"`), "401 レスポンスが重複している")
}

func TestListRecommendUsers_InternalError(t *testing.T) {
	uc := &mockRecommendUsecase{
		listRecommend: func(_ context.Context, _ uint) ([]*usecase.RecommendUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(&mockPickupUsecase{}, &mockNewUserUsecase{}, uc), http.MethodGet, "/users/recommend", "", 1)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListRecommendUsers_Empty(t *testing.T) {
	uc := &mockRecommendUsecase{
		listRecommend: func(_ context.Context, _ uint) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(&mockPickupUsecase{}, &mockNewUserUsecase{}, uc), http.MethodGet, "/users/recommend", "", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "users")
}

func TestListRecommendUsers_NoProfileLeak(t *testing.T) {
	uc := &mockRecommendUsecase{
		listRecommend: func(_ context.Context, _ uint) ([]*usecase.RecommendUserResult, error) {
			return []*usecase.RecommendUserResult{
				{UserID: 1, Name: "田中"},
			}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(&mockPickupUsecase{}, &mockNewUserUsecase{}, uc), http.MethodGet, "/users/recommend", "", 1)

	body := rec.Body.String()
	assert.NotContains(t, body, "email")
	assert.NotContains(t, body, "password")
	assert.NotContains(t, body, "aliases")
}
