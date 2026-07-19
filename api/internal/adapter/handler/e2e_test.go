package handler_test

// E2E tests: real TCP server (httptest.Server) + real WsHub + mock usecases
// DB は不要。usecase 層をモックして handler + router 層を通しで検証する。

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/domain/model"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mock: SendMessageUsecase ---

type mockSendMessageUsecase struct {
	fn func(ctx context.Context, roomID, senderUserID uint, content string) (*model.Message, error)
}

func (m *mockSendMessageUsecase) SendMessage(ctx context.Context, roomID, senderUserID uint, content string) (*model.Message, error) {
	if m.fn != nil {
		return m.fn(ctx, roomID, senderUserID, content)
	}
	return &model.Message{ID: 1, RoomID: roomID, SenderUserID: senderUserID, Content: content, CreatedAt: time.Now()}, nil
}

// --- mock: InterestUsecase (マッチ通知 E2E 用) ---

type mockInterestUsecase struct {
	fn func(ctx context.Context, from, to uint) (*usecase.SendInterestResult, error)
}

func (m *mockInterestUsecase) SendInterest(ctx context.Context, from, to uint) (*usecase.SendInterestResult, error) {
	if m.fn != nil {
		return m.fn(ctx, from, to)
	}
	return &usecase.SendInterestResult{Matched: false}, nil
}

// --- helper: build real echo server for E2E ---

func newE2EServer(hub *handler.WsHub, msgUc usecase.SendMessageUsecase) *httptest.Server {
	h := &handler.Handler{
		AuthHandler:    handler.NewAuthHandler(&mockAuthUsecase{}),
		HealthHandler:  handler.NewHealthHandler(),
		MessageHandler: handler.NewMessageHandler(msgUc, hub),
	}
	return httptest.NewServer(router.NewRouter(h, hub))
}

func newE2EServerWithInterest(hub *handler.WsHub, intUc usecase.InterestUsecase) *httptest.Server {
	h := &handler.Handler{
		AuthHandler:    handler.NewAuthHandler(&mockAuthUsecase{}),
		HealthHandler:  handler.NewHealthHandler(),
		MessageHandler: handler.NewMessageHandler(&mockSendMessageUsecase{}, hub),
		UsersHandler: handler.NewUsersHandler(
			&mockPickupUsecase{},
			&mockNewUserUsecase{},
			&mockRecommendUsecase{},
			nil, nil,
			intUc,
			nil,
			hub,
		),
	}
	return httptest.NewServer(router.NewRouter(h, hub))
}

func e2ePost(t *testing.T, srv *httptest.Server, path, body, auth string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, srv.URL+path, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

// --- E2E: POST /messages → WS /ws/rooms/:roomId に push ---

func TestE2E_SendMessage_BroadcastsToWSRoom(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newE2EServer(hub, &mockSendMessageUsecase{
		fn: func(_ context.Context, roomID, senderUserID uint, content string) (*model.Message, error) {
			return &model.Message{ID: 42, RoomID: roomID, SenderUserID: senderUserID, Content: content, CreatedAt: time.Now()}, nil
		},
	})
	defer srv.Close()

	wsConn := dialWS(t, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws/rooms/1")
	defer wsConn.Close()
	waitConnRegistered(func() bool { return hub.RoomConnCount(1) == 1 })

	resp := e2ePost(t, srv, "/messages", `{"roomId":1,"content":"E2Eテスト"}`, testBearerToken(10))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	msg := readWSJSON(t, wsConn)
	assert.Equal(t, "new_message", msg["type"])
	assert.Equal(t, float64(1), msg["roomId"])
	payload := msg["message"].(map[string]any)
	assert.Equal(t, "E2Eテスト", payload["content"])
	assert.Equal(t, float64(10), payload["senderUserId"])
}

// --- E2E: POST /messages → 別ルームのWS接続には届かない ---

func TestE2E_SendMessage_DoesNotLeakToOtherRoom(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newE2EServer(hub, &mockSendMessageUsecase{})
	defer srv.Close()

	wsRoom1 := dialWS(t, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws/rooms/1")
	defer wsRoom1.Close()
	wsRoom2 := dialWS(t, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws/rooms/2")
	defer wsRoom2.Close()
	waitConnRegistered(func() bool { return hub.RoomConnCount(1) == 1 && hub.RoomConnCount(2) == 1 })

	resp := e2ePost(t, srv, "/messages", `{"roomId":1,"content":"room1only"}`, testBearerToken(10))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// ルーム1は受信
	msg := readWSJSON(t, wsRoom1)
	assert.Equal(t, "new_message", msg["type"])

	// ルーム2はタイムアウト（受信しない）
	wsRoom2.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, _, err := wsRoom2.ReadMessage()
	assert.Error(t, err)
}

// --- E2E: 複数の WS クライアントが同じルームにいる場合は全員に届く ---

func TestE2E_SendMessage_BroadcastsToAllClientsInRoom(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newE2EServer(hub, &mockSendMessageUsecase{})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/rooms/5"
	ws1 := dialWS(t, wsURL)
	defer ws1.Close()
	ws2 := dialWS(t, wsURL)
	defer ws2.Close()
	waitConnRegistered(func() bool { return hub.RoomConnCount(5) == 2 })

	resp := e2ePost(t, srv, "/messages", `{"roomId":5,"content":"全員へ"}`, testBearerToken(1))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	msg1 := readWSJSON(t, ws1)
	msg2 := readWSJSON(t, ws2)
	assert.Equal(t, "new_message", msg1["type"])
	assert.Equal(t, "new_message", msg2["type"])
}

// --- E2E: POST /messages 認証なし → 401 / WS には何も届かない ---

func TestE2E_SendMessage_Unauthorized_NoWsBroadcast(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newE2EServer(hub, &mockSendMessageUsecase{})
	defer srv.Close()

	wsConn := dialWS(t, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws/rooms/1")
	defer wsConn.Close()
	waitConnRegistered(func() bool { return hub.RoomConnCount(1) == 1 })

	// 認証ヘッダーなし
	resp := e2ePost(t, srv, "/messages", `{"roomId":1,"content":"unauthorized"}`, "")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// WS には何も届かない
	wsConn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, _, err := wsConn.ReadMessage()
	assert.Error(t, err)
}

// --- E2E: SendInterest がマッチした場合 WS /ws/rooms に通知が届く ---

func TestE2E_SendInterest_Match_NotifiesRoomListWS(t *testing.T) {
	hub := handler.NewWsHub()
	roomID := uint(99)
	srv := newE2EServerWithInterest(hub, &mockInterestUsecase{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			return &usecase.SendInterestResult{Matched: true, RoomID: &roomID}, nil
		},
	})
	defer srv.Close()

	wsConn := dialWS(t, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws/rooms")
	defer wsConn.Close()
	waitConnRegistered(func() bool { return hub.RoomListConnCount() == 1 })

	resp := e2ePost(t, srv, "/users/2/interests", "", testBearerToken(1))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	msg := readWSJSON(t, wsConn)
	assert.Equal(t, "rooms_updated", msg["type"])
}

// --- E2E: SendInterest がマッチしない場合 WS /ws/rooms には通知されない ---

func TestE2E_SendInterest_NoMatch_NoWsNotification(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newE2EServerWithInterest(hub, &mockInterestUsecase{
		fn: func(_ context.Context, _, _ uint) (*usecase.SendInterestResult, error) {
			return &usecase.SendInterestResult{Matched: false}, nil
		},
	})
	defer srv.Close()

	wsConn := dialWS(t, "ws"+strings.TrimPrefix(srv.URL, "http")+"/ws/rooms")
	defer wsConn.Close()
	waitConnRegistered(func() bool { return hub.RoomListConnCount() == 1 })

	resp := e2ePost(t, srv, "/users/2/interests", "", testBearerToken(1))
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	wsConn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, _, err := wsConn.ReadMessage()
	assert.Error(t, err, "マッチしていないので rooms_updated は届かないはず")
}
