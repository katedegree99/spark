package model

import "time"

type UserInterest struct {
	FromUserID uint `gorm:"primaryKey"`
	ToUserID   uint `gorm:"primaryKey"`
	CreatedAt  time.Time
}
