package model

import "time"

type UserPickupHistory struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	UserID      uint `gorm:"not null;uniqueIndex:uk_user_shown"`
	ShownUserID uint `gorm:"not null;uniqueIndex:uk_user_shown"`
	CreatedAt   time.Time
}
