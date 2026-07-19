package model

import "time"

type Room struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"not null"`
}

type RoomUser struct {
	RoomID    uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
}

type Message struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	RoomID       uint      `gorm:"not null;index"`
	SenderUserID uint      `gorm:"not null"`
	Content      string    `gorm:"not null;type:text"`
	CreatedAt    time.Time `gorm:"not null"`
}
