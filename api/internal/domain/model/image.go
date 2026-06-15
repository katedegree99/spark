package model

import "time"

type Image struct {
	ID        uint      `gorm:"primaryKey"`
	Directory string    `gorm:"not null"`
	URL       string    `gorm:"not null"`
	CreatedAt time.Time
}
