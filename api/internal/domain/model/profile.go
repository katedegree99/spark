package model

import "time"

type Profile struct {
	UserID      uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"`
	IconImageID *uint
	Bio         *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
