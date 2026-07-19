package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/domain/model"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEchoMessage(uc usecase.SendMessageUsecase) *echo.Echo {
	hub := handler.NewWsHub()
	h := &handler.Handler{
		AuthHandler:    handler.NewAuthHandler(&mockAuthUsecase{}),
		HealthHandler:  handler.NewHealthHandler(),
		MessageHandler: handler.NewMessageHandler(uc, hub),
	}
	return router.NewRouter(h, hub)
}

func TestSendMessage_Success(t *testing.T) {
	now := time.Now()
	uc := &mockSendMessageUsecase{
		fn: func(_ context.Context, roomID, senderUserID uint, content string) (*model.Message, error) {
			return &model.Message{ID: 1, RoomID: roomID, SenderUserID: senderUserID, Content: content, CreatedAt: now}, nil
		},
	}
	rec := doWithAuth(
		newTestEchoMessage(uc),
		http.MethodPost, "/messages",
		`{"roomId":1,"content":"こんにちは"}`,
		42,
	)
	assert.Equal(t, http.StatusCreated, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, "こんにちは")
	require.Contains(t, body, `"roomId"`)
	require.Contains(t, body, `"senderUserId"`)
	require.Contains(t, body, `"isMine"`)
}

func TestSendMessage_Unauthorized(t *testing.T) {
	uc := &mockSendMessageUsecase{}
	rec := do(newTestEchoMessage(uc), http.MethodPost, "/messages", `{"roomId":1,"content":"test"}`)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSendMessage_NotRoomMember_Returns403(t *testing.T) {
	uc := &mockSendMessageUsecase{
		fn: func(_ context.Context, _, _ uint, _ string) (*model.Message, error) {
			return nil, usecase.ErrNotRoomMember
		},
	}
	rec := doWithAuth(
		newTestEchoMessage(uc),
		http.MethodPost, "/messages",
		`{"roomId":99,"content":"test"}`,
		1,
	)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "FORBIDDEN")
}

func TestSendMessage_InternalError(t *testing.T) {
	uc := &mockSendMessageUsecase{
		fn: func(_ context.Context, _, _ uint, _ string) (*model.Message, error) {
			return nil, errors.New("db error")
		},
	}
	rec := doWithAuth(
		newTestEchoMessage(uc),
		http.MethodPost, "/messages",
		`{"roomId":1,"content":"test"}`,
		1,
	)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSendMessage_IsMineAlwaysTrue(t *testing.T) {
	uc := &mockSendMessageUsecase{
		fn: func(_ context.Context, roomID, senderUserID uint, content string) (*model.Message, error) {
			return &model.Message{ID: 5, RoomID: roomID, SenderUserID: senderUserID, Content: content, CreatedAt: time.Now()}, nil
		},
	}
	rec := doWithAuth(
		newTestEchoMessage(uc),
		http.MethodPost, "/messages",
		`{"roomId":1,"content":"hello"}`,
		99,
	)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"isMine":true`)
}
