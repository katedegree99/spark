package model

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Email          string `gorm:"uniqueIndex;not null"`
	HashedPassword string `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
