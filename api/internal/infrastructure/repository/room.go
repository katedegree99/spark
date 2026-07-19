package repository

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/model"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *roomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, userID1, userID2 uint) (*model.Room, error) {
	room := &model.Room{}
	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		return nil, err
	}
	members := []model.RoomUser{
		{RoomID: room.ID, UserID: userID1},
		{RoomID: room.ID, UserID: userID2},
	}
	if err := r.db.WithContext(ctx).Create(&members).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (r *roomRepository) FindByUsers(ctx context.Context, userID1, userID2 uint) (*model.Room, error) {
	var roomID uint
	err := r.db.WithContext(ctx).
		Table("room_users ru1").
		Select("ru1.room_id").
		Joins("JOIN room_users ru2 ON ru1.room_id = ru2.room_id AND ru2.user_id = ?", userID2).
		Where("ru1.user_id = ?", userID1).
		Limit(1).
		Scan(&roomID).Error
	if err != nil {
		return nil, err
	}
	if roomID == 0 {
		return nil, nil
	}
	var room model.Room
	if err := r.db.WithContext(ctx).First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) IsMember(ctx context.Context, roomID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.RoomUser{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	return count > 0, err
}
