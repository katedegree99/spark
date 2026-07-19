package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsRoomsEvent struct {
	Type string `json:"type"`
}

type wsMessageEvent struct {
	Type    string           `json:"type"`
	RoomID  uint             `json:"roomId"`
	Message WsMessagePayload `json:"message"`
}

type WsMessagePayload struct {
	ID           uint   `json:"id"`
	SenderUserID uint   `json:"senderUserId"`
	Content      string `json:"content"`
	CreatedAt    string `json:"createdAt"`
}

// WsHub はルームと接続中クライアントを管理する
type WsHub struct {
	// roomID -> connections
	rooms map[uint]map[*websocket.Conn]struct{}
	// rooms list connections
	roomListConns map[*websocket.Conn]struct{}
}

func NewWsHub() *WsHub {
	return &WsHub{
		rooms:         make(map[uint]map[*websocket.Conn]struct{}),
		roomListConns: make(map[*websocket.Conn]struct{}),
	}
}

// NotifyRoomsUpdated はルーム一覧の更新をルーム一覧 WS 接続に通知する
func (h *WsHub) NotifyRoomsUpdated() {
	event, _ := json.Marshal(wsRoomsEvent{Type: "rooms_updated"})
	for conn := range h.roomListConns {
		conn.WriteMessage(websocket.TextMessage, event)
	}
}

// BroadcastMessage は指定ルームの接続にメッセージを push する
func (h *WsHub) BroadcastMessage(roomID uint, payload WsMessagePayload) {
	event, _ := json.Marshal(wsMessageEvent{
		Type:    "new_message",
		RoomID:  roomID,
		Message: payload,
	})
	for conn := range h.rooms[roomID] {
		conn.WriteMessage(websocket.TextMessage, event)
	}
}

// RoomListConnCount はテスト用にルーム一覧接続数を返す
func (h *WsHub) RoomListConnCount() int {
	return len(h.roomListConns)
}

// RoomConnCount はテスト用に指定ルームの接続数を返す
func (h *WsHub) RoomConnCount(roomID uint) int {
	return len(h.rooms[roomID])
}

// HandleRoomsWS は WS /ws/rooms を処理する
func (h *WsHub) HandleRoomsWS(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		delete(h.roomListConns, conn)
		conn.Close()
	}()
	h.roomListConns[conn] = struct{}{}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
	return nil
}

// HandleRoomWS は WS /ws/rooms/:roomId を処理する
func (h *WsHub) HandleRoomWS(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	var roomID uint
	if _, err := fmt.Sscan(c.Param("roomId"), &roomID); err != nil {
		conn.Close()
		return err
	}

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*websocket.Conn]struct{})
	}
	h.rooms[roomID][conn] = struct{}{}

	defer func() {
		delete(h.rooms[roomID], conn)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
	return nil
}
