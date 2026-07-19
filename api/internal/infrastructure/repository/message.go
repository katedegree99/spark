package repository

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/model"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, roomID, senderUserID uint, content string) (*model.Message, error) {
	msg := &model.Message{
		RoomID:       roomID,
		SenderUserID: senderUserID,
		Content:      content,
	}
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}
