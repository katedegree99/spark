package handler

import (
	"context"
	"errors"

	authmw "github.com/katedegree99/spark/api/internal/adapter/middleware"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/katedegree99/spark/api/pkg/generated"
)

type MessageHandler struct {
	sendMessageUsecase usecase.SendMessageUsecase
	hub                *WsHub
}

func NewMessageHandler(sendMessageUsecase usecase.SendMessageUsecase, hub *WsHub) *MessageHandler {
	return &MessageHandler{sendMessageUsecase: sendMessageUsecase, hub: hub}
}

func (h *MessageHandler) SendMessage(ctx context.Context, req generated.SendMessageRequestObject) (generated.SendMessageResponseObject, error) {
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.SendMessage401JSONResponse{UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{Code: "UNAUTHORIZED", Message: "unauthorized"}}, nil
	}

	msg, err := h.sendMessageUsecase.SendMessage(ctx, uint(req.Body.RoomId), userID, req.Body.Content)
	if err != nil {
		if errors.Is(err, usecase.ErrNotRoomMember) {
			return generated.SendMessage403JSONResponse{
				Code:    "FORBIDDEN",
				Message: "not a room member",
			}, nil
		}
		return nil, err
	}

	h.hub.BroadcastMessage(msg.RoomID, WsMessagePayload{
		ID:           msg.ID,
		SenderUserID: msg.SenderUserID,
		Content:      msg.Content,
		CreatedAt:    msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})

	isMine := true
	return generated.SendMessage201JSONResponse{
		Id:           int(msg.ID),
		RoomId:       int(msg.RoomID),
		SenderUserId: int(msg.SenderUserID),
		Content:      msg.Content,
		IsMine:       isMine,
		CreatedAt:    msg.CreatedAt,
	}, nil
}
