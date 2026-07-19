package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newWsTestServer(hub *handler.WsHub) *httptest.Server {
	e := echo.New()
	e.GET("/ws/rooms", hub.HandleRoomsWS)
	e.GET("/ws/rooms/:roomId", hub.HandleRoomWS)
	return httptest.NewServer(e)
}

func wsURL(srv *httptest.Server, path string) string {
	return "ws" + strings.TrimPrefix(srv.URL, "http") + path
}

func dialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	return conn
}

func readWSJSON(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)
	var result map[string]any
	require.NoError(t, json.Unmarshal(msg, &result))
	return result
}

// 接続が hub に登録されるまで待つ
func waitConnRegistered(check func() bool) bool {
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if check() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// --- WS /ws/rooms ---

func TestWsHub_NotifyRoomsUpdated_SingleClient(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn := dialWS(t, wsURL(srv, "/ws/rooms"))
	defer conn.Close()

	waitConnRegistered(func() bool { return hub.RoomListConnCount() == 1 })

	hub.NotifyRoomsUpdated()

	msg := readWSJSON(t, conn)
	assert.Equal(t, "rooms_updated", msg["type"])
}

func TestWsHub_NotifyRoomsUpdated_MultipleClients(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn1 := dialWS(t, wsURL(srv, "/ws/rooms"))
	defer conn1.Close()
	conn2 := dialWS(t, wsURL(srv, "/ws/rooms"))
	defer conn2.Close()

	waitConnRegistered(func() bool { return hub.RoomListConnCount() == 2 })

	hub.NotifyRoomsUpdated()

	msg1 := readWSJSON(t, conn1)
	msg2 := readWSJSON(t, conn2)
	assert.Equal(t, "rooms_updated", msg1["type"])
	assert.Equal(t, "rooms_updated", msg2["type"])
}

func TestWsHub_HandleRoomsWS_DisconnectCleansUp(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn := dialWS(t, wsURL(srv, "/ws/rooms"))
	waitConnRegistered(func() bool { return hub.RoomListConnCount() == 1 })

	conn.Close()

	waitConnRegistered(func() bool { return hub.RoomListConnCount() == 0 })
	assert.Equal(t, 0, hub.RoomListConnCount())
}

// --- WS /ws/rooms/:roomId ---

func TestWsHub_BroadcastMessage_ReceivesInCorrectRoom(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn := dialWS(t, wsURL(srv, "/ws/rooms/1"))
	defer conn.Close()

	waitConnRegistered(func() bool { return hub.RoomConnCount(1) == 1 })

	hub.BroadcastMessage(1, handler.WsMessagePayload{
		ID:           42,
		SenderUserID: 10,
		Content:      "こんにちは",
		CreatedAt:    "2026-07-20T00:00:00Z",
	})

	msg := readWSJSON(t, conn)
	assert.Equal(t, "new_message", msg["type"])
	assert.Equal(t, float64(1), msg["roomId"])
	payload := msg["message"].(map[string]any)
	assert.Equal(t, "こんにちは", payload["content"])
	assert.Equal(t, float64(10), payload["senderUserId"])
}

func TestWsHub_BroadcastMessage_OtherRoomNotAffected(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn1 := dialWS(t, wsURL(srv, "/ws/rooms/1"))
	defer conn1.Close()
	conn2 := dialWS(t, wsURL(srv, "/ws/rooms/2"))
	defer conn2.Close()

	waitConnRegistered(func() bool { return hub.RoomConnCount(1) == 1 && hub.RoomConnCount(2) == 1 })

	hub.BroadcastMessage(1, handler.WsMessagePayload{Content: "ルーム1のメッセージ"})

	// ルーム1は受信できる
	msg := readWSJSON(t, conn1)
	assert.Equal(t, "new_message", msg["type"])

	// ルーム2は受信しない（タイムアウトになる）
	conn2.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, _, err := conn2.ReadMessage()
	assert.Error(t, err, "ルーム2は受信すべきでない")
}

func TestWsHub_BroadcastMessage_MultipleClientsInSameRoom(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn1 := dialWS(t, wsURL(srv, "/ws/rooms/5"))
	defer conn1.Close()
	conn2 := dialWS(t, wsURL(srv, "/ws/rooms/5"))
	defer conn2.Close()

	waitConnRegistered(func() bool { return hub.RoomConnCount(5) == 2 })

	hub.BroadcastMessage(5, handler.WsMessagePayload{Content: "全員に届く"})

	msg1 := readWSJSON(t, conn1)
	msg2 := readWSJSON(t, conn2)
	assert.Equal(t, "new_message", msg1["type"])
	assert.Equal(t, "new_message", msg2["type"])
}

func TestWsHub_HandleRoomWS_DisconnectCleansUp(t *testing.T) {
	hub := handler.NewWsHub()
	srv := newWsTestServer(hub)
	defer srv.Close()

	conn := dialWS(t, wsURL(srv, "/ws/rooms/3"))
	waitConnRegistered(func() bool { return hub.RoomConnCount(3) == 1 })

	conn.Close()

	waitConnRegistered(func() bool { return hub.RoomConnCount(3) == 0 })
	assert.Equal(t, 0, hub.RoomConnCount(3))
}
