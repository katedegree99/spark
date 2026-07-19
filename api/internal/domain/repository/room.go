package repository

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/model"
)

type RoomRepository interface {
	Create(ctx context.Context, userID1, userID2 uint) (*model.Room, error)
	FindByUsers(ctx context.Context, userID1, userID2 uint) (*model.Room, error)
	IsMember(ctx context.Context, roomID, userID uint) (bool, error)
}
