package usecase

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/model"
	"github.com/katedegree99/spark/api/internal/domain/repository"
)

var ErrNotRoomMember = errors.New("not a room member")

type SendMessageUsecase interface {
	SendMessage(ctx context.Context, roomID, senderUserID uint, content string) (*model.Message, error)
}

type sendMessageUsecase struct {
	roomRepo    repository.RoomRepository
	messageRepo repository.MessageRepository
}

func NewSendMessageUsecase(roomRepo repository.RoomRepository, messageRepo repository.MessageRepository) SendMessageUsecase {
	return &sendMessageUsecase{roomRepo: roomRepo, messageRepo: messageRepo}
}

func (u *sendMessageUsecase) SendMessage(ctx context.Context, roomID, senderUserID uint, content string) (*model.Message, error) {
	isMember, err := u.roomRepo.IsMember(ctx, roomID, senderUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotRoomMember
	}
	return u.messageRepo.Create(ctx, roomID, senderUserID, content)
}
