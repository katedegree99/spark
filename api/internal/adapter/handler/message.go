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
}

func NewMessageHandler(sendMessageUsecase usecase.SendMessageUsecase) *MessageHandler {
	return &MessageHandler{sendMessageUsecase: sendMessageUsecase}
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
