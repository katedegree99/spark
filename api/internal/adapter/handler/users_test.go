package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/katedegree/spark/api/internal/adapter/handler"
	"github.com/katedegree/spark/api/internal/adapter/router"
	"github.com/katedegree/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPickupUsecase struct {
	listPickup func(ctx context.Context, userID uint) ([]*usecase.PickupUserResult, error)
}

func (m *mockPickupUsecase) ListPickup(ctx context.Context, userID uint) ([]*usecase.PickupUserResult, error) {
	if m.listPickup != nil {
		return m.listPickup(ctx, userID)
	}
	return nil, nil
}

func newTestEchoWithUsers(uc usecase.PickupUsecase) *echo.Echo {
	h := &handler.Handler{
		UsersHandler: handler.NewUsersHandler(uc),
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
	rec := doWithAuth(newTestEchoWithUsers(uc), http.MethodGet, "/users/pickup", "", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "堺 理人")
	require.Contains(t, body, "matchedTags")
	require.Contains(t, body, "unmatchedTags")
}

func TestListPickupUsers_Unauthorized(t *testing.T) {
	uc := &mockPickupUsecase{}
	rec := do(newTestEchoWithUsers(uc), http.MethodGet, "/users/pickup", "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestListPickupUsers_InternalError(t *testing.T) {
	uc := &mockPickupUsecase{
		listPickup: func(_ context.Context, _ uint) ([]*usecase.PickupUserResult, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(uc), http.MethodGet, "/users/pickup", "", 1)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListPickupUsers_Empty(t *testing.T) {
	uc := &mockPickupUsecase{
		listPickup: func(_ context.Context, _ uint) ([]*usecase.PickupUserResult, error) {
			return []*usecase.PickupUserResult{}, nil
		},
	}
	rec := doWithAuth(newTestEchoWithUsers(uc), http.MethodGet, "/users/pickup", "", 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "users")
}
