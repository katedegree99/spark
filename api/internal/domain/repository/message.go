package repository

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/model"
)

type MessageRepository interface {
	Create(ctx context.Context, roomID, senderUserID uint, content string) (*model.Message, error)
}
